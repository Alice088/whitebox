package prompts

import "time"

type Case struct {
	Name   string
	Input  string
	Prompt string
}

type Metrics struct {
	EventCalls int
	Steps      int
	ToolCalls  int
	Errors     int
	Duration   time.Duration
}

type Result struct {
	Name    string
	Output  string
	Error   error
	Metrics Metrics
}
