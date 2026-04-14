package tui

import (
	"context"
	"fmt"
	"os"
	"time"
	syscontext "whitebox/internal/core/context"
	"whitebox/internal/core/llm"
	"whitebox/internal/core/tracer"
	"whitebox/internal/core/tui/status"
	"whitebox/internal/core/wrapfuse"

	tea "charm.land/bubbletea/v2"
	"github.com/henomis/langfuse-go/model"
	"github.com/rs/zerolog"
)

type Chat struct {
	LLM          llm.LLM
	Context      *syscontext.Context
	Logger       zerolog.Logger
	statusEngine *status.StatusEngine
	Tracer       tracer.Tracer
}

func New(llm llm.LLM, ctx *syscontext.Context, session syscontext.Sessions, tracer tracer.Tracer, logger zerolog.Logger) Chat {
	return Chat{
		Context:      ctx,
		LLM:          llm,
		Logger:       logger,
		Session:      session,
		statusEngine: status.NewStatusEngine(),
		Tracer:       tracer,
	}
}

func (chat *Chat) Run(ctx context.Context) {
	m := initialModel(chat, chat.Session.ID)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
		chat.Logger.Error().Err(err).Msg("failed to start tui")
	}
}

func (chat *Chat) ask(ctx context.Context, input string) (string, error) {
	lf := chat.Wrapfuse.NewClient(ctx)
	defer func() {
		chat.Wrapfuse.Flush(lf, ctx)
	}()

	trace, err := chat.Wrapfuse.Trace(lf, &model.Trace{
		Name:      "whitebox-chat",
		Input:     input,
		Timestamp: new(time.Now()),
	})
	if err != nil {
		chat.Logger.Error().Err(err).Msg("failed to create langfuse tracee")
		return chat.LLM.Ask(input, chat.Context.Prompt())
	}

	traceID := ""
	if chat.Wrapfuse.Enabled {
		traceID = trace.ID
	}

	generation, err := chat.Wrapfuse.Generation(lf, &model.Generation{
		Model:   chat.LLM.Model(),
		Name:    "llm-call",
		TraceID: traceID,
		Input: []model.M{
			{"role": "system", "content": chat.Context.Prompt()},
			{"role": "user", "content": input},
		},
	}, nil)
	if err != nil {
		chat.Logger.Error().Err(err).Msg("failed to create langfuse generation")
		return chat.LLM.Ask(input, chat.Context.Prompt())
	}

	generationID := ""
	if chat.Wrapfuse.Enabled {
		generationID = generation.ID
	}

	output, err := chat.LLM.Ask(input, chat.Context.Prompt())
	if err != nil {
		_, gErr := chat.Wrapfuse.GenerationEnd(lf, &model.Generation{
			ID:     generationID,
			Output: model.M{"error": err.Error()},
		})
		if gErr != nil {
			chat.Logger.Error().Err(err).Msg("failed to end langfuse generation")
		}
		return "", err
	}

	generation.Output = model.M{"completion": output}
	generation.Usage = model.Usage{
		Input:  int(chat.LLM.EstimateTokens(input)),
		Output: int(chat.LLM.EstimateTokens(output)),
		Total:  int(chat.LLM.EstimateTokens(input + chat.Context.Prompt() + output)),
	}

	_, err = chat.Wrapfuse.GenerationEnd(lf, generation)
	if err != nil {
		chat.Logger.Error().Err(err).Msg("failed to end langfuse generation")
	}

	_, err = chat.Wrapfuse.Trace(lf, &model.Trace{
		ID:     traceID,
		Output: output,
	})
	if err != nil {
		chat.Logger.Error().Err(err).Msg("failed to update langfuse trace")
	}

	return output, nil
}
