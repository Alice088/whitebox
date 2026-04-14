package tools

import (
	"fmt"
)

func Execute(call ToolCall) (string, error) {
	var tool Tool
	var ok bool

	if tool, ok = Tools[call.Tool]; !ok {
		return "", fmt.Errorf("unknown tool: %s", call.Tool)
	}

	return tool(call.Arguments)
}
