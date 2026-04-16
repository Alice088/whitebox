package secure

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"whitebox/internal/paths"
)

func Path(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}

	if filepath.IsAbs(path) {
		return "", fmt.Errorf("absolute paths not allowed")
	}

	var root string

	if strings.HasPrefix(path, "memory/") {
		root = paths.MemoriesDir
		path = strings.TrimPrefix(path, "memory/")
	} else {
		root = paths.WorkspaceDir
	}

	fullPath := filepath.Join(root, path)
	cleanPath := filepath.Clean(fullPath)

	rel, err := filepath.Rel(root, cleanPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("access denied")
	}

	realPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		realPath = cleanPath
	}

	rel, err = filepath.Rel(root, realPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("access denied (symlink)")
	}

	return realPath, nil
}
