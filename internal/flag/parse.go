package flag

import (
	"errors"
	"flag"
	"strings"
	"whitebox/internal/factory"
)

type Input struct {
	Model      string
	Provider   factory.ProviderOpts
	SessionID  string
	MaxHistory int
}

func ParseFlags() (Input, error) {
	model := flag.String("model", "", "model name")
	provider := flag.String("provider", "local", "provider: api | local")
	providerName := flag.String("provider_name", "", "provider name like: deepseek, llamacpp")
	sessionID := flag.String("session", "", "session ID for persistent chat")
	maxHistory := flag.Int("max-history", 10, "maximum messages to keep in history")

	flag.Parse()

	if *model == "" {
		return Input{}, errors.New("model required")
	}

	normalizedProvider := strings.TrimSpace(strings.ToLower(*provider))
	if normalizedProvider != string(factory.APIProvider) && normalizedProvider != string(factory.LocalProvider) {
		return Input{}, errors.New("provider must be 'api' or 'local'")
	}

	normalizedProviderName := strings.TrimSpace(strings.ToLower(*providerName))
	if normalizedProvider == string(factory.APIProvider) && normalizedProviderName == "" {
		return Input{}, errors.New("provider_name required for api provider")
	}

	return Input{
		Model: *model,
		Provider: factory.ProviderOpts{
			ProviderType: factory.ToProvider(normalizedProvider),
			Name:         normalizedProviderName,
		},
		SessionID:  *sessionID,
		MaxHistory: *maxHistory,
	}, nil
}
