package llm

type LLM interface {
	Ask(prompt string, id string) (string, error)
	EstimateTokens(string) float64
	Model() string
}
