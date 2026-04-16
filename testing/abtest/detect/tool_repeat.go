package detect

import (
	"maps"
	"whitebox/internal/core/tools"
)

func ToolRepeat(toolCalls []tools.ToolCall) bool {
	var prevTool tools.ToolCall
	for _, tool := range toolCalls {
		if tool.Tool == prevTool.Tool && maps.Equal(tool.Arguments, prevTool.Arguments) {
			return true
		}
	}
	return false
}
