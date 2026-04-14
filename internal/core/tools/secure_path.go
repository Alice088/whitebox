package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func securePath(baseDir, path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}

	// запрещаем абсолютные пути
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("absolute paths not allowed")
	}

	// собираем полный путь
	fullPath := filepath.Join(baseDir, path)
	cleanPath := filepath.Clean(fullPath)

	// проверка выхода через ../
	rel, err := filepath.Rel(baseDir, cleanPath)
	if err != nil {
		return "", fmt.Errorf("invalid path")
	}
	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("access denied")
	}

	// резолвим symlink
	realPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil {
		// если файл еще не существует (write_file) — это нормально
		if !os.IsNotExist(err) {
			return "", err
		}
		realPath = cleanPath
	}

	// повторная проверка после symlink
	rel, err = filepath.Rel(baseDir, realPath)
	if err != nil {
		return "", fmt.Errorf("invalid path")
	}
	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("access denied (symlink)")
	}

	return realPath, nil
}
