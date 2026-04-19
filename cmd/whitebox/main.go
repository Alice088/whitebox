package main

import (
	"encoding/json"
	"fmt"
	"whitebox/internal/broker/task"
	"whitebox/internal/core"
	syscontext "whitebox/internal/core/context"
	"whitebox/internal/core/embedded_prompts"
	"whitebox/internal/factory"
	"whitebox/internal/langfuse"
	"whitebox/internal/providers"
	"whitebox/pkg/cfg"
	"whitebox/pkg/logging"
	"whitebox/pkg/meta"
	"whitebox/pkg/prepare"

	"github.com/mailru/easyjson"
	"github.com/nats-io/nats.go"
)

func init() {
	err := prepare.EnsureWhitebox()
	if err != nil {
		panic(fmt.Sprintf("Failed to init .whitebox-%s", meta.AgentName))
	}
}

func main() {
	logger := logging.MustLogger()
	config := cfg.MustConfig(logger)

	nc, err := nats.Connect(config.NATS.URL)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	_, err = nc.Subscribe("task.created.*", func(msg *nats.Msg) {
		var t task.Created
		if err := easyjson.Unmarshal(msg.Data, &t); err != nil {
			return
		}

		if t.Type != config.Agent.Type {
			return
		}

		session := syscontext.NewSession(t.TaskID, config.Session)
		if err := session.CreateSessionDir(); err != nil {
			return
		}
		session.MustLoadMessages(&logger)

		systemContext := syscontext.New(session, embedded_prompts.OutputProtocolV1())
		if err := systemContext.Collect(); err != nil {
			return
		}

		llm, err := factory.LLM(config.LLM.Provide, providers.InitOpts{
			Model:  config.LLM.Model,
			ApiKey: config.LLM.ApiKey,
		})
		if err != nil {
			return
		}

		if config.Observability.LangFuse.Enabled {
			llm = &langfuse.LLMWrapper{LLM: llm}
		}

		engine := core.Engine{
			LLM:     llm,
			Context: systemContext,
			CallChain: core.CallChain{
				Max: config.CallChain.Max,
			},
		}

		answer, err := engine.Run(t.Payload.Msg, func(event core.Event) {
			raw, _ := json.Marshal(event)
			nc.Publish("task.logs."+t.TaskID, raw)
		})

		if err != nil {
			nc.Publish("task.error."+t.TaskID, []byte(err.Error()))
			return
		}

		res := task.Result{
			TaskID: t.TaskID,
			Result: answer,
		}

		raw, _ := easyjson.Marshal(res)
		nc.Publish("task.result."+t.TaskID, raw)
	})

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	select {}
}
