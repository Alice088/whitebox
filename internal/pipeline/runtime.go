package pipeline

import (
	syscontext "whitebox/internal/core/context"
	llmcore "whitebox/internal/core/llm"
)

type Runtime struct {
	LLM     llmcore.LLM
	Context syscontext.Context
}
