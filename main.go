package main

import (
	"whitebox/internal/core"
	syscontext "whitebox/internal/core/context"
	"whitebox/internal/factory"
	"whitebox/internal/flag"
	"whitebox/internal/providers"
	"whitebox/internal/tui"
	"whitebox/pkg/cfg"
	"whitebox/pkg/logging"
	"whitebox/pkg/wbinit"
)

func init() {
	baseDir, err := wbinit.EnsureWhitebox()
	if err != nil {
		panic("Failed to init .whitebox")
	}
	syscontext.BaseDir = baseDir + "/context"
}

func main() {

	logger := logging.MustLogger()
	input := flag.MustInput(logger)
	config := cfg.MustConfig(logger)

	session := syscontext.NewSession(input.SessionID, config.Session)
	err := session.CreateSessionDir()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
	session.MustLoadMessages(&logger)

	systemContext := syscontext.Context{
		Sessions: session,
	}

	err = systemContext.Collect()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	llm, err := factory.LLM(input.Provider, providers.InitOpts{
		Model:  input.Model,
		ApiKey: config.LLM.ApiKey,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init LLM")
	}

	engine := core.Engine{
		LLM:     llm,
		Context: &systemContext,
		CallChain: core.CallChain{
			Max: config.CallChain.Max,
		},
	}

	chat := tui.New(engine, input.Debug)
	chat.Run()
}
