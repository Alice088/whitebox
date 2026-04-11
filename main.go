package main

import (
	"coreclaw/internal/config"
	"coreclaw/internal/llm/deepseek"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
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

	var cfg config.Config
	err = env.Parse(&cfg)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	llm := deepseek.New(deepseek.Opts{
		Model:  "deepseek-reasoner",
		ApiKey: cfg.LLM.ApiKey,
	})
	answer, err := llm.Ask(prompt)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to ask LLM")
	}

	fmt.Printf("> %s\n", answer)
}
