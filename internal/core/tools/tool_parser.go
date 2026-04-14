package tools

import "encoding/json"

type ToolCall struct {
	Tool      string    `json:"tool"`
	Arguments Arguments `json:"arguments,omitzero"`
}

type Arguments struct {
	Path    string `json:"path"`
	Content string `json:"content"`
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
