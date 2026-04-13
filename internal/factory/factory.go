package factory

import (
	"fmt"
	"strings"
	"sync"
	llm2 "whitebox/internal/core/llm"
	"whitebox/internal/providers"
)

type Provider string

const (
	APIProvider   Provider = "api"
	LocalProvider Provider = "local"
)

const (
	defaultLocalProvider = "llamacpp"
)

type ProviderOpts struct {
	Name         string
	ProviderType Provider
}

type Constructor func(initOpts providers.InitOpts) llm2.LLM

var (
	mu             sync.RWMutex
	apiProviders   = make(map[string]Constructor)
	localProviders = make(map[string]Constructor)
)

func RegisterAPI(name string, constructor Constructor) {
	registerProvider(apiProviders, name, constructor)
}

func RegisterLocal(name string, constructor Constructor) {
	registerProvider(localProviders, name, constructor)
}

func LLM(providerOpts ProviderOpts, initOpts providers.InitOpts) (llm2.LLM, error) {
	registerProviders()

	providerName := normalizeProviderName(providerOpts)
	constructor, err := resolveConstructor(providerOpts.ProviderType, providerName)
	if err != nil {
		return nil, err
	}

	return constructor(initOpts), nil
}

func resolveConstructor(providerType Provider, providerName string) (Constructor, error) {
	mu.RLock()
	defer mu.RUnlock()

	switch providerType {
	case LocalProvider:
		constructor, ok := localProviders[providerName]
		if !ok {
			return nil, fmt.Errorf("unknown local provider %q", providerName)
		}
		return constructor, nil
	case APIProvider:
		constructor, ok := apiProviders[providerName]
		if !ok {
			return nil, fmt.Errorf("unknown api provider %q", providerName)
		}
		return constructor, nil
	default:
		return nil, fmt.Errorf("unknown provider type %q", providerType)
	}
}

func registerProvider(target map[string]Constructor, name string, constructor Constructor) {
	if constructor == nil {
		panic("provider constructor is nil")
	}

	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		panic("provider name is empty")
	}

	mu.Lock()
	defer mu.Unlock()

	target[name] = constructor
}

func normalizeProviderName(opts ProviderOpts) string {
	name := strings.TrimSpace(strings.ToLower(opts.Name))
	if opts.ProviderType == LocalProvider && name == "" {
		return defaultLocalProvider
	}
	return name
}

func ToProvider(provider string) Provider {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case string(LocalProvider):
		return LocalProvider
	case string(APIProvider):
		return APIProvider
	default:
		return APIProvider
	}
}
