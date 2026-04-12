package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"whitebox/internal/config"
	xcontext "whitebox/internal/context"
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

	flags, err := flag.ParseFlags()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	err = godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	l := langfuse.New(context.Background())

	var cfg config.Config
	err = env.Parse(&cfg)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	systemContext := xcontext.Context{}
	err = systemContext.Collect(xcontext.CollectOpts{
		ToolsPath:  "./context/tools",
		MindPath:   "./context/mind",
		MemoryPath: "./context/memory",
		//MessagesPath: "nope",
		SkillsPath: "./context/skills",
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load context")
	}

	llm, err := factory.New(flags.Provider, xllm.InitOpts{
		Model:  flags.Model,
		ApiKey: cfg.LLM.ApiKey,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to init LLM")
	}

	t, err := l.Trace(&model.Trace{
		Name:      "coreclaw-request",
		Input:     flags.Msg,
		Timestamp: new(time.Now()),
	})
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	systemPrompt := systemContext.Prompt()
	g, err := l.Generation(&model.Generation{
		Model:   llm.Model(),
		Name:    "llm-call",
		TraceID: t.ID,
		Input: []model.M{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": flags.Msg},
		},
	}, nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create generation")
	}

	answer, err := llm.Ask(flags.Msg, systemPrompt)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to ask LLM")
	}

	g.Output = model.M{"completion": answer}
	g.Usage = model.Usage{
		Input:  int(llm.EstimateTokens(flags.Msg)),
		Output: int(llm.EstimateTokens(answer)),
		Total:  int(llm.EstimateTokens(answer + flags.Msg + systemPrompt)),
	}
	_, err = l.GenerationEnd(g)
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
