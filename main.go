package main

import (
	"context"
	"whitebox/internal/core/chat"
	syscontext "whitebox/internal/core/context"
	"whitebox/internal/factory"
	"whitebox/internal/flag"
	"whitebox/internal/providers"
	"whitebox/pkg/cfg"
	"whitebox/pkg/logging"
)

func main() {
	logger := logging.MustLogger()
	input := flag.MustInput(logger)
	config := cfg.MustConfig(logger)

	session := syscontext.NewSession(input.SessionID, input.MaxHistory, logger)
	systemContext := syscontext.NewMustDefault(session.MustLoadMessages(), logger)

	llm, err := factory.LLM(input.Provider, providers.InitOpts{
		Model:  input.Model,
		ApiKey: config.LLM.ApiKey,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init LLM")
	}

	tui := chat.New(llm, systemContext, session, logger)
	tui.Run(context.Background())
}
