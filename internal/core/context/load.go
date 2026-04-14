package context

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func load(dir string) ([]Item, error) {
	var items []Item

	if dir == "" {
		return items, nil
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return items, nil
	}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to context file: %w", err)
		}

		items = append(items, Item{
			Source:  path,
			Content: strings.TrimSpace(string(data)),
		})

		return nil
	})

	if err != nil {
		return items, fmt.Errorf("failed to walk on dir: %w", err)
	}

	return items, nil
}
