package main

import (
	"context"
	"fmt"
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
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logRotator := &lumberjack.Logger{
		Filename:   "./logs/whitebox.log",
		MaxSize:    10,
		MaxBackups: 2,
		MaxAge:     28,
		Compress:   true,
	}

	logger := zerolog.New(logRotator).With().Timestamp().Logger()

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

	runner := &pipeline.Runner{Logger: &logger}
	runner.Use("langfuse_start", pipeline.LangfuseStart(lf, "whitebox-request"))
	runner.Use("build_prompt", pipeline.BuildPrompt())
	runner.Use("ask_llm", pipeline.AskLLM())
	runner.Use("langfuse_end", pipeline.LangfuseEnd(lf))
	runner.Use("logging", pipeline.Logging(logger))

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
