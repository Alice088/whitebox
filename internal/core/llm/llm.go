package llm

type LLM interface {
	Ask(prompt, systemPrompt string) (string, error)
	EstimateTokens(string) float64
	Model() string
}
