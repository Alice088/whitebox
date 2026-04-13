package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const maxFileSize = 1_000_000 // ~1MB

func ReadFile(path string) (string, error) {
	baseDir := "/home/gosha/.whitebox/workspace"

	if path == "" {
		return "", fmt.Errorf("empty path")
	}

	if filepath.IsAbs(path) {
		return "", fmt.Errorf("absolute paths not allowed")
	}

	fullPath := filepath.Join(baseDir, path)
	cleanPath := filepath.Clean(fullPath)

	rel, err := filepath.Rel(baseDir, cleanPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("access denied")
	}

	realPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil {
		return "", err
	}

	rel, err = filepath.Rel(baseDir, realPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("access denied (symlink)")
	}

	info, err := os.Stat(realPath)
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

	f, err := os.Open(realPath)
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
