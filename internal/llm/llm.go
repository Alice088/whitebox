package llm

type LLM interface {
	Ask(prompt string, id string) (string, error)
	Model() string
}
