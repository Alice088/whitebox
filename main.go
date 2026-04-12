package main

import (
	"context"
	"fmt"
	"os"
	cfg "whitebox/internal/config"
	syscontext "whitebox/internal/context"
	"whitebox/internal/factory"
	"whitebox/internal/flag"
	"whitebox/internal/pipeline"
	xllm "whitebox/internal/providers"

	"github.com/caarlos0/env/v11"
	"github.com/henomis/langfuse-go"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout)

	input, err := flag.ParseFlags()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	err = godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	var config cfg.Config
	err = env.Parse(&config)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	systemContext, err := syscontext.NewDefault()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load context")
	}

	llm, err := factory.LLM(input.Provider, xllm.InitOpts{
		Model:  input.Model,
		ApiKey: config.LLM.ApiKey,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init LLM")
	}

	lf := langfuse.New(context.Background())

	runner := &pipeline.Runner{}
	runner.Read(pipeline.Logging(logger))
	runner.Write(pipeline.LangfuseStart(lf, "whitebox-request"))
	runner.Write(pipeline.BuildPrompt())
	runner.Write(pipeline.AskLLM())
	runner.Write(pipeline.LangfuseEnd(lf))

	state := &pipeline.State{
		Input:   input.Msg,
		LLM:     llm,
		Context: systemContext,
	}

	err = runner.Run(context.Background(), state)
	if err != nil {
		logger.Fatal().Err(err).Msg("pipeline failed")
	}

	lf.Flush(context.Background())
	fmt.Printf("> %s\n", state.Output)
}
