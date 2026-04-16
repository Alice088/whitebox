package abtest

import (
	"time"
	"whitebox/internal/core/tools"
)

type Case struct {
	Name   string
	Input  string
	Prompt string
}

type Metrics struct {
	EventCalls        int
	Steps             int
	ToolsCalls        map[string]int
	ToolsCallsHistory []tools.ToolCall
	Errors            int
	Duration          time.Duration
}

type Result struct {
	Name    string
	Output  string
	Error   error
	Metrics Metrics
}
