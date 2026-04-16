package main

import (
	"os"
	"whitebox/internal/core"
	"whitebox/internal/factory"
	"whitebox/internal/providers"
	"whitebox/pkg/cfg"
	"whitebox/testing/abtest"
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

	cases := []abtest.Case{
		{
			Name:   "normal",
			Input:  "create file test.txt",
			Prompt: plusBase(load("./testing/abtest/prompts/files/normal.md")),
		},
		//{
		//	Name:   "loop_force",
		//	Input:  "create file test.txt",
		//	Prompt: plusBase(load("./testing/abtest/prompts/files/loop_force.md")),
		//},
		{
			Name:   "extra_use",
			Input:  "create file test.txt",
			Prompt: plusBase(load("./testing/abtest/prompts/files/extra_use.md")),
		},
		//{
		//	Name:   "analyst",
		//	Input:  "create file test.txt",
		//	Prompt: plusBase(load("./testing/abtest/prompts/files/analyst.md")),
		//},
	}

	logger.Info().Msg("Run butch")
	results := runner.RunBatch(cases)

	prompts.PrintReport(results)
}

func plusBase(prompt string) string {
	return load("./testing/abtest/prompts/files/base.md") + prompt
}

func load(path string) string {
	b, _ := os.ReadFile(path)
	return string(b)
}
