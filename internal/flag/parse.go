package flag

import (
	"errors"
	"flag"
	"strings"
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
	providerName := flag.String("provider_name", "", "provider name like: deepseek, llamacpp")

	flag.Parse()

	if *model == "" {
		return Config{}, errors.New("model required")
	}

	if *msg == "" {
		return Config{}, errors.New("msg is empty")
	}

	normalizedProvider := strings.TrimSpace(strings.ToLower(*provider))
	if normalizedProvider != string(factory.APIProvider) && normalizedProvider != string(factory.LocalProvider) {
		return Config{}, errors.New("provider must be 'api' or 'local'")
	}

	normalizedProviderName := strings.TrimSpace(strings.ToLower(*providerName))
	if normalizedProvider == string(factory.APIProvider) && normalizedProviderName == "" {
		return Config{}, errors.New("provider_name required for api provider")
	}

	return Config{
		Model: *model,
		Provider: factory.ProviderOpts{
			ProviderType: factory.ToProvider(normalizedProvider),
			Name:         normalizedProviderName,
		},
		Msg: *msg,
	}, nil
}
