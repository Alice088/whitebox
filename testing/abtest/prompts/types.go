package prompts

type Case struct {
	Name   string
	Input  string
	Prompt string
}

type Result struct {
	Name   string
	Output string
	Error  error
}
