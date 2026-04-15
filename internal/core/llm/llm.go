package llm

import "whitebox/internal/core/tools"

type LLM interface {
	Ask(prompt, systemPrompt string) (string, error)
	Tool(call tools.ToolCall) (string, error)
	EstimateTokens(string) float64
	Model() string
}
