package abtest

import (
	"fmt"
	"whitebox/internal/core/tools"
	"whitebox/pkg/messages"
)

func StringToolCallsHistory(history []tools.ToolCall) string {
	if len(history) == 0 {
		return "no calls"
	}

	str := "[ "

	for i, tool := range history {

		if i == len(history)-1 {
			str += fmt.Sprintf("%s(%+v) ] -> len(%d)", tool.Tool, messages.StringArgs(tool.Arguments), len(history))
		} else {
			str += fmt.Sprintf("%s(%+v) -> ", tool.Tool, messages.StringArgs(tool.Arguments))
		}
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
