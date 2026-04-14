package tools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"whitebox/internal/core/tools/secure"
	"whitebox/pkg/maps"
)

// WriteFile - fields: path, content
func WriteFile(arguments map[string]string) (string, error) {
	if !maps.Exists(arguments, "path") {
		return "", errors.New("path field required")
	}

	if !maps.Exists(arguments, "content") {
		return "", errors.New("content field required")
	}

	path, err := secure.Path(arguments["path"])
	if err != nil {
		return "", fmt.Errorf("unsecure path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(path, []byte(arguments["content"]), 0644); err != nil {
		return "", err
	}

	return fmt.Sprintf("file written: %s", path), nil
}
