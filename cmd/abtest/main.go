package main

import (
	"os"
	"whitebox/internal/core"
	"whitebox/internal/factory"
	"whitebox/internal/providers"
	"whitebox/pkg/cfg"
	"whitebox/testing/abtest/prompts"

	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = os.Stdout
	}))
	config := cfg.MustConfig(logger)

	llm, err := factory.LLM(factory.ProviderOpts{
		"deepseek",
		factory.APIProvider,
	}, providers.InitOpts{
		Model:  "deepseek-reasoner",
		ApiKey: config.LLM.ApiKey,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create llm")
	}

	engine := core.Engine{
		LLM: llm,
		CallChain: core.CallChain{
			Max: config.CallChain.Max,
		},
	}

	runner := prompts.Runner{
		Engine: &engine,
		Logger: &logger,
	}

	cases := []prompts.Case{
		{
			Name:   "nice prompt: v1",
			Input:  "create file test.txt",
			Prompt: load("./testing/abtest/prompts/files/example_v1.md"),
		},
		{
			Name:   "bad prompt: v2",
			Input:  "create file test.txt",
			Prompt: load("./testing/abtest/prompts/files/example_v2.md"),
		},
	}

	logger.Info().Msg("Run butch")
	results := runner.RunBatch(cases)

	prompts.PrintReport(results)
}

func load(path string) string {
	b, _ := os.ReadFile(path)
	return string(b)
}
