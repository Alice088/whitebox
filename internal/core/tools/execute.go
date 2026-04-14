package tools

func Execute(call ToolCall) (string, error) {
	switch call.Tool {
	case "read_file":
		return ReadFile(call.Arguments.Path)
	}
	return "", nil
}
