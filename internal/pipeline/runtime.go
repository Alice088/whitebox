package pipeline

import (
	syscontext "whitebox/internal/context"
	llmcore "whitebox/internal/core/llm"
)

type Runtime struct {
	LLM     llmcore.LLM
	Context syscontext.Context
}
