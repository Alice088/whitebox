package detect

import (
	"whitebox/internal/core/tools"
)

func ToolLoop(toolCalls []tools.ToolCall) bool {
	var prevTool tools.ToolCall
	for _, tool := range toolCalls {
		if tool.Tool == prevTool.Tool {
			return true
		}
		prevTool = tool
	}
	return false
}
