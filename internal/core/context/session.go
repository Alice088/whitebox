package context

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"whitebox/internal/config"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Sessions struct {
	Messages    []Message
	ID          string
	MaxMessages int
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewSession(ID string, config config.Session) Sessions {
	s := Sessions{
		MaxMessages: config.MaxMessages,
	}

	if ID != "" {
		s.ID = ID
	} else {
		s.ID = uuid.NewString()
	}

	return s
}

func (s *Sessions) MustLoadMessages(logger *zerolog.Logger) {
	data, err := os.ReadFile(s.Path())
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		logger.Fatal().Err(err).Msg("Failed to load session")
	}

	err = json.Unmarshal(data, &s.Messages)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to unmarshal messages")
	}

	logger.Info().
		Str("session_id", s.ID).
		Int("loaded_messages", len(s.Messages)).
		Msg("session loaded")
}

func (s *Sessions) SaveSession(msgs []Message) error {
	dir := filepath.Dir(s.Path())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(msgs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.Path(), data, 0644)
}

func (s *Sessions) CreateSessionDir() error {
	if err := os.MkdirAll(SessionsDir, 0755); err != nil {
		return err
	}
	return nil
}

func (s *Sessions) Path() string {
	return fmt.Sprintf("%s/%s.json", SessionsDir, s.ID)
}
