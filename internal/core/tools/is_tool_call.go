package tools

import "encoding/json"

func IsToolCall(call string) (ToolCall, bool) {
	var tc ToolCall
	if err := json.Unmarshal([]byte(call), &tc); err != nil {
		return ToolCall{}, false
	}
	if tc.Tool == "" {
		return ToolCall{}, false
	}
	return tc, true
}
