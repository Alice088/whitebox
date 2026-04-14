package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteFile(path, content string) (string, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}

	return fmt.Sprintf("file written: %s", path), nil
}
