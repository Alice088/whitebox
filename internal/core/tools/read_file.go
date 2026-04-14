package tools

import (
	"errors"
	"fmt"
	"io"
	"os"
	"whitebox/internal/core/tools/secure"
	"whitebox/pkg/maps"
)

const maxFileSize = 1_000_000 // ~1MB

// ReadFile - fields: path
func ReadFile(arguments map[string]string) (string, error) {
	if !maps.Exists(arguments, "path") {
		return "", errors.New("path field required")
	}

	path, err := secure.Path(arguments["path"])
	if err != nil {
		return "", fmt.Errorf("unsecure path: %w", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return "", fmt.Errorf("cannot read directory")
	}

	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("not a regular file")
	}

	if info.Size() > maxFileSize {
		return "", fmt.Errorf("file too large")
	}

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	limited := io.LimitReader(f, maxFileSize+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return "", err
	}

	if int64(len(data)) > maxFileSize {
		return "", fmt.Errorf("file too large")
	}

	return string(data), nil
}
