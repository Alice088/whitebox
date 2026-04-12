package flag

import (
	"errors"
	"flag"
	"whitebox/internal/llm/factory"
)

type Config struct {
	Msg      string
	Model    string
	Provider factory.ProviderOpts
}

func ParseFlags() (Config, error) {
	msg := flag.String("msg", "", "message to llm")
	model := flag.String("model", "", "model name")
	provider := flag.String("provider", "local", "provider: api | local")
	providerName := flag.String("provider_name", "", "provider name like: deepseek, openai")

	flag.Parse()

	if *model == "" {
		return Config{}, errors.New("model required")
	}

	if *msg == "" {
		return Config{}, errors.New("msg is empty")
	}

	if *providerName == "" {
		return Config{}, errors.New("provider_name required")
	}

	if *provider != "api" && *provider != "local" {
		return Config{}, errors.New("provider must be 'api' or 'local'")
	}

	return Config{
		Model: *model,
		Provider: factory.ProviderOpts{
			ProviderType: factory.ToProvider(*provider),
			Name:         *providerName,
		},
		Msg: *msg,
	}, nil
}
