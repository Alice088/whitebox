package main

import (
	"context"
	"fmt"
	"os"
	"time"
	cfg "whitebox/internal/config"
	syscontext "whitebox/internal/context"
	"whitebox/internal/factory"
	"whitebox/internal/flag"
	xllm "whitebox/internal/providers"

	"github.com/caarlos0/env/v11"
	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
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

	l := langfuse.New(context.Background())

	var config cfg.Config
	err = env.Parse(&config)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	systemContext, err := syscontext.NewDefault()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load context")
	}

	llm, err := factory.LLM(input.Provider, xllm.InitOpts{
		Model:  input.Model,
		ApiKey: config.LLM.ApiKey,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to init LLM")
	}

	t, err := l.Trace(&model.Trace{
		Name:      "coreclaw-request",
		Input:     input.Msg,
		Timestamp: new(time.Now()),
	})
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	generation, err := l.Generation(&model.Generation{
		Model:   llm.Model(),
		Name:    "llm-call",
		TraceID: t.ID,
		Input: []model.M{
			{"role": "system", "content": systemContext.Prompt()},
			{"role": "user", "content": input.Msg},
		},
	}, nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create generation")
	}

	answer, err := llm.Ask(input.Msg, systemContext.Prompt())
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to ask LLM")
	}

	generation.Output = model.M{"completion": answer}
	generation.Usage = model.Usage{
		Input:  int(llm.EstimateTokens(input.Msg)),
		Output: int(llm.EstimateTokens(answer)),
		Total:  int(llm.EstimateTokens(answer + input.Msg + systemContext.Prompt())),
	}
	_, err = l.GenerationEnd(generation)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to generation_end")
	}

	_, err = l.Trace(&model.Trace{
		ID:     t.ID,
		Output: answer,
	})
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	l.Flush(context.TODO())
	fmt.Printf("> %s\n", answer)
}
