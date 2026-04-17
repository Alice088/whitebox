package factory

import (
	"fmt"
	"strings"
	"sync"

	llm2 "whitebox/internal/core/llm"
	"whitebox/internal/providers"
)

type Constructor func(initOpts providers.InitOpts) llm2.LLM

var (
	mu       sync.RWMutex
	registry = map[string]Constructor{}
)

func Register(name string, constructor Constructor) {
	if constructor == nil {
		panic("provider constructor is nil")
	}

	name = normalize(name)
	if name == "" {
		panic("provider name is empty")
	}

	mu.Lock()
	defer mu.Unlock()

	registry[name] = constructor
}

func LLM(name string, initOpts providers.InitOpts) (llm2.LLM, error) {
	registerProviders()

	name = normalize(name)

	if name == "" {
		name = "llamacpp"
	}

	mu.RLock()
	constructor, ok := registry[name]
	mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown provider %q", name)
	}

	return constructor(initOpts), nil
}

func normalize(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}
