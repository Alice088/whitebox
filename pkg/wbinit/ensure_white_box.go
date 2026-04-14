package wbinit

import (
	"os"
	"path/filepath"
)

func EnsureWhitebox() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	base := filepath.Join(home, ".whitebox")

	dirs := []string{
		filepath.Join(base, "context", "minds"),
		filepath.Join(base, "context", "skills"),
		filepath.Join(base, "context", "tools"),
		filepath.Join(base, "context", "memories"),
		filepath.Join(base, "context", "sessions"),
		filepath.Join(base, "workspace"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return "", err
		}
	}

	return base, nil
}
