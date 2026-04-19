package main

import (
	"fmt"
	"whitebox/internal/config"
	"whitebox/internal/core"
	syscontext "whitebox/internal/core/context"
	"whitebox/internal/core/embedded_prompts"
	"whitebox/internal/factory"
	"whitebox/internal/flags"
	"whitebox/internal/langfuse"
	"whitebox/internal/providers"
	"whitebox/pkg/cfg"
	"whitebox/pkg/logging"
	"whitebox/pkg/prepare"
)

func init() {
	err := prepare.EnsureWhitebox()
	if err != nil {
		panic(fmt.Sprintf("Failed to init .whitebox-%s", config.AgentName))
	}
}

func main() {
	logger := logging.MustLogger()
	flag := flags.MustInput(logger)
	config := cfg.MustConfig(logger)

	session := syscontext.NewSession(flag.SessionID, config.Session)
	err := session.CreateSessionDir()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
	session.MustLoadMessages(&logger)

	systemContext := syscontext.New(session, embedded_prompts.OutputProtocolV1())
	err = systemContext.Collect()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	llm, err := factory.LLM(config.LLM.Provide, providers.InitOpts{
		Model:  config.LLM.Model,
		ApiKey: config.LLM.ApiKey,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init LLM")
	}

	if config.Observability.LangFuse.Enabled {
		llm = &langfuse.LLMWrapper{
			LLM: llm,
		}
	}

	engine := core.Engine{
		LLM:     llm,
		Context: systemContext,
		CallChain: core.CallChain{
			Max: config.CallChain.Max,
		},
	}

	answer, err := engine.Run(flag.Msg, func(event core.Event) {
		fmt.Printf("%+v\n", event)
	})
	if err != nil {
		panic("run err: " + err.Error())
	}
	fmt.Println(answer)
	return
}
