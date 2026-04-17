package factory

import (
	"whitebox/internal/providers/deepseek"
	"whitebox/internal/providers/llamacpp"
)

func registerProviders() {
	Register("deepseek", deepseek.New)
	Register("llamacpp", llamacpp.New)
}
