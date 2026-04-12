package factory

import (
	"coreclaw/internal/llm"
	"coreclaw/internal/llm/deepseek"
	"coreclaw/internal/llm/llamacpp"
	"errors"
)

type Provider string

const (
	APIProvider   Provider = "api"
	LocalProvider Provider = "local"
)

type ProviderOpts struct {
	Name         string
	ProviderType Provider
}

func New(providerOpts ProviderOpts, initOpts llm.InitOpts) (llm.LLM, error) {
	switch providerOpts.ProviderType {
	case LocalProvider:
		return newLocal(initOpts), nil
	case APIProvider:
		return newAPI(providerOpts, initOpts), nil
	default:
		return nil, errors.New("unknow provider")
	}
}

func newLocal(opts llm.InitOpts) llm.LLM {
	return llamacpp.New(opts)
}

func newAPI(providerOpts ProviderOpts, opts llm.InitOpts) llm.LLM {
	switch providerOpts.Name {
	case "deepseek":
		return deepseek.New(opts)
	default:
		return deepseek.New(opts)
	}
}

func ToProvider(provider string) Provider {
	switch provider {
	case "local":
		return LocalProvider
	case "API":
		return APIProvider
	default:
		return APIProvider
	}
}
