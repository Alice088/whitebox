package prepare

import (
	"os"
	"path/filepath"
	"whitebox/internal/paths"
)

func EnsureWhitebox() error {
	dirs := []string{
		filepath.Join(paths.BaseDir, "context", "minds"),
		filepath.Join(paths.BaseDir, "context", "skills"),
		filepath.Join(paths.BaseDir, "context", "tools"),
		filepath.Join(paths.BaseDir, "context", "memories"),
		filepath.Join(paths.BaseDir, "context", "sessions"),
		filepath.Join(paths.BaseDir, "commands"),
		filepath.Join(paths.BaseDir, "workspace"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	return nil
}
