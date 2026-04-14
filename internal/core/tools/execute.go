package tools

import (
	"fmt"
	"whitebox/internal/core/context"
)

func Execute(call ToolCall) (string, error) {
	if call.Arguments.Path == "" {
		return "", fmt.Errorf("read_file: path required")
	}

	realPath, err := securePath(context.WorkspaceDir, call.Arguments.Path)
	if err != nil {
		return "", err
	}

	switch call.Tool {
	case "read_file":
		return ReadFile(realPath)

	case "write_file":
		return WriteFile(realPath, call.Arguments.Content)

	default:
		return "", fmt.Errorf("unknown tool: %s", call.Tool)
	}
}
