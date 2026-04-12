package llm

type LLM interface {
	Ask(prompt string, systemPrompt string) (string, error)
	EstimateTokens(string) float64
	Model() string
}
