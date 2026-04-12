package context

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func NewSessionID() string {
	return uuid.NewString()
}

func LoadSession(path string) ([]Message, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Message{}, nil
		}
		return nil, err
	}

	var msgs []Message
	err = json.Unmarshal(data, &msgs)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func SaveSession(path string, msgs []Message) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(msgs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
