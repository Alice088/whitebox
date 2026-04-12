package main

import (
	"context"
	"coreclaw/internal/config"
	xcontext "coreclaw/internal/context"
	"coreclaw/internal/llm/llamacpp"
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

	ctx := xcontext.Context{}
	err = ctx.Collect(xcontext.CollectOpts{
		ToolsPath:  "./context/tools",
		MindPath:   "./context/mind",
		MemoryPath: "./context/memory",
		//MessagesPath: "nope",
		SkillsPath: "./context/skills",
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load context")
	}

	llm := llamacpp.New(llamacpp.Opts{
		Model: "Gemopus-4-E4B-it-Preview-Q6_K.gguf",
		//ApiKey:   cfg.LLM.ApiKey,
		LangFuse: l,
		Logger:   &logger,
		Context:  ctx,
	})

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
