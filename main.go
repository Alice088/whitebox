package main

import (
	"context"
	"coreclaw/internal/config"
	"coreclaw/internal/llm/deepseek"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(0)
	}
	prompt := os.Args[1]

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout)

	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	l := langfuse.New(context.Background())

	var cfg config.Config
	err = env.Parse(&cfg)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	llm := deepseek.New(deepseek.Opts{
		Model:    "deepseek-reasoner",
		ApiKey:   cfg.LLM.ApiKey,
		LangFuse: l,
		Logger:   &logger,
	}, "отвечай кратко и по делу")

	t, err := l.Trace(&model.Trace{
		Name:      "coreclaw-request",
		Input:     prompt,
		Timestamp: new(time.Now()),
	})

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	answer, err := llm.Ask(prompt, t.ID)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to ask LLM")
	}

	_, _ = l.Trace(&model.Trace{
		ID:     t.ID,
		Output: answer,
	})

	l.Flush(context.TODO())
	fmt.Printf("> %s\n", answer)
}
