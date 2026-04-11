package llm

type LLM interface {
	Ask(string) (string, error)
	Model() string
}
