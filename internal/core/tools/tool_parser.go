package tools

import "encoding/json"

type ToolCall struct {
	Tool      string    `json:"tool"`
	Arguments Arguments `json:"arguments"`
}

type Arguments struct {
	Path string `json:"path"`
}

func TryParseToolCall(s string) (ToolCall, bool) {
	var tc ToolCall
	if err := json.Unmarshal([]byte(s), &tc); err != nil {
		return ToolCall{}, false
	}
	if tc.Tool == "" {
		return ToolCall{}, false
	}
	return tc, true
}
