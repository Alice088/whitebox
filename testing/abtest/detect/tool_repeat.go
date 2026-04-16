package detect

import (
	"fmt"
	"whitebox/internal/core/tools"
)

func ToolRepeat(toolCalls []tools.ToolCall) int {
	seen := make(map[string]int)
	loops := 0

	for _, t := range toolCalls {
		key := t.Tool + fmt.Sprintf("%v", t.Arguments)

		seen[key]++

		if seen[key] > 1 {
			loops++
		}
	}

	return loops
}
