package tools

import (
	"fmt"
	"io"
	"os"
)

const maxFileSize = 1_000_000 // ~1MB

func ReadFile(path string) (string, error) {
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
