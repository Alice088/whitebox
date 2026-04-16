package abtest

import (
	"fmt"
	"strings"
	"whitebox/internal/core/tools"
)

func StringToolCallsHistory(history []tools.ToolCall) string {
	if len(history) == 0 {
		return "no calls"
	}

	str := "[ "

	for i, tool := range history {

		if i == len(history)-1 {
			str += fmt.Sprintf("%s(%+v) ] -> len(%d)", tool.Tool, stringArgs(tool.Arguments), len(history))
		} else {
			str += fmt.Sprintf("%s(%+v) -> ", tool.Tool, stringArgs(tool.Arguments))
		}
	}

	return str
}

func stringArgs(args map[string]string) string {
	if len(args) == 0 {
		return ""
	}

	var parts []string

	for k, v := range args {
		parts = append(parts, fmt.Sprintf("%s:%s", k, v))
	}

	return strings.Join(parts, ", ")
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
