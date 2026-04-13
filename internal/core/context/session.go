package context

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Session struct {
	Logger      zerolog.Logger
	ID          string
	Path        string
	MaxMessages int
}

func NewSession(ID string, max int, logger zerolog.Logger) Session {
	s := Session{
		Logger:      logger,
		MaxMessages: max,
	}

	if ID != "" {
		s.ID = ID
	} else {
		s.ID = uuid.NewString()
	}

	s.CreateSessionDir()
	return s
}

func (s *Session) MustLoadMessages() []Message {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Message{}
		}
		s.Logger.Fatal().Err(err).Msg("Failed to load session")
	}

	var msgs []Message
	err = json.Unmarshal(data, &msgs)
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("Failed to unmarshal messages")
	}

	s.Logger.Info().
		Str("session_id", s.ID).
		Int("loaded_messages", len(msgs)).
		Msg("session loaded")

	return msgs
}

func (s *Session) SaveSession(msgs []Message) error {
	dir := filepath.Dir(s.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(msgs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.Path, data, 0644)
}

func (s *Session) CreateSessionDir(path ...string) {
	sessionsDir := "context/sessions"
	if len(path) != 0 {
		sessionsDir = path[0]
	}

	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		s.Logger.Fatal().Err(err).Msg("failed to create sessions directory")
	}

	s.Path = filepath.Join(sessionsDir, s.ID+".json")
	return
}
