package llm

import (
	"coreclaw/internal/context"

	"github.com/henomis/langfuse-go"
	"github.com/rs/zerolog"
)

type InitOpts struct {
	ApiKey   string
	BaseURL  string
	Model    string
	LangFuse *langfuse.Langfuse
	Logger   *zerolog.Logger
	Context  context.Context
}
