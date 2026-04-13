package factory

import (
	"whitebox/internal/providers/deepseek"
	"whitebox/internal/providers/llamacpp"
)

func registerProviders() {
	RegisterAPI("deepseek", deepseek.New)
	RegisterLocal(defaultLocalProvider, llamacpp.New)
}
