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

	command := arguments["command"]
	if command == "" {
		return "", errors.New("empty command")
	}

	if err := secure.Command(command); err != nil {
		return "", fmt.Errorf("unsecure command: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-lc", command)
	cmd.Dir = paths.WorkspaceDir
	sys.SetSysProcAttr(cmd)

	output, err := cmd.CombinedOutput()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return string(output), errors.New("command timeout")
	}

	if err != nil {
		return string(output), fmt.Errorf("execution error: %w", err)
	}

	out := strings.TrimSpace(string(output))

	if len(out) == 0 {
		out = "OK"
	}

	return out, nil
}
