package testing

import (
	"fmt"
	"whitebox/internal/core/tools"
)

func StringToolCallsHistory(history []tools.ToolCall) string {
	if len(history) == 0 {
		return "no calls"
	}

	str := "[ "

	for i, tool := range history {
		if i == len(history)-1 {
			str += fmt.Sprintf("%s ] -> (%d)", tool.Tool, len(history))
		}

		str += fmt.Sprintf("%s ", tool.Tool)
	}

	return str
}

func StringToolCalls(calls map[string]int) string {
	if len(calls) == 0 {
		return "no calls"
	}

	str := ""

	for name, count := range calls {
		str += fmt.Sprintf("%s:%d ", name, count)
	}

	return str
}
