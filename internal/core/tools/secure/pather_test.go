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

	// запрещаем абсолютные пути
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("absolute paths not allowed")
	}

	// собираем полный путь
	fullPath := filepath.Join(paths.WorkspaceDir, path)
	cleanPath := filepath.Clean(fullPath)

	// проверка выхода через ../
	rel, err := filepath.Rel(paths.WorkspaceDir, cleanPath)
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
	rel, err = filepath.Rel(paths.WorkspaceDir, realPath)
	if err != nil {
		return "", fmt.Errorf("invalid path")
	}
	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("access denied (symlink)")
	}

	return realPath, nil
}
