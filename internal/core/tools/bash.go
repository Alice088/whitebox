package tools

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
	"whitebox/internal/core/tools/secure"
	"whitebox/internal/paths"
	"whitebox/pkg/maps"
	"whitebox/pkg/sys"
)

// Bash - fields: command
func Bash(arguments map[string]string) (string, error) {
	if !maps.Exists(arguments, "command") {
		return "", errors.New("command field required")
	}

	if err := secure.Command(arguments["command"]); err != nil {
		return "", fmt.Errorf("unsecure command: %w", err)
	}

	parts := strings.Fields(arguments["command"])
	if len(parts) == 0 {
		return "", errors.New("invalid command")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Dir = paths.WorkspaceDir
	sys.SetSysProcAttr(cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("execution error: %w", err)
	}

	if len(output) == 0 {
		output = []byte("OK")
	}

	return string(output), nil
}
