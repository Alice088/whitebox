package tui

import (
	"context"
	"fmt"
	"os"
	"time"
	syscontext "whitebox/internal/core/context"
	"whitebox/internal/core/llm"
	"whitebox/internal/core/tui/status"

	tea "charm.land/bubbletea/v2"
	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
	"github.com/rs/zerolog"
)

type Chat struct {
	LLM          llm.LLM
	Context      *syscontext.Context
	Logger       zerolog.Logger
	Session      syscontext.Session
	statusEngine *status.StatusEngine
}

func New(llm llm.LLM, ctx *syscontext.Context, session syscontext.Session, logger zerolog.Logger) Chat {
	return Chat{
		Context:      ctx,
		LLM:          llm,
		Logger:       logger,
		Session:      session,
		statusEngine: status.NewStatusEngine(),
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
	lf := langfuse.New(ctx)
	defer lf.Flush(ctx)

	trace, err := lf.Trace(&model.Trace{
		Name:      "whitebox-chat",
		Input:     input,
		Timestamp: new(time.Now()),
	})
	if err != nil {
		chat.Logger.Error().Err(err).Msg("failed to create langfuse tracee")
		return chat.LLM.Ask(input, chat.Context.Prompt())
	}

	g, err := lf.Generation(&model.Generation{
		Model:   chat.LLM.Model(),
		Name:    "llm-call",
		TraceID: trace.ID,
		Input: []model.M{
			{"role": "system", "content": chat.Context.Prompt()},
			{"role": "user", "content": input},
		},
	}, nil)
	if err != nil {
		chat.Logger.Error().Err(err).Msg("failed to create langfuse generation")
		return chat.LLM.Ask(input, chat.Context.Prompt())
	}

	output, err := chat.LLM.Ask(input, chat.Context.Prompt())
	if err != nil {
		_, gErr := lf.GenerationEnd(&model.Generation{
			ID:     g.ID,
			Output: model.M{"error": err.Error()},
		})
		if gErr != nil {
			chat.Logger.Error().Err(err).Msg("failed to end langfuse generation")
		}
		return "", err
	}

	g.Output = model.M{"completion": output}
	g.Usage = model.Usage{
		Input:  int(chat.LLM.EstimateTokens(input)),
		Output: int(chat.LLM.EstimateTokens(output)),
		Total:  int(chat.LLM.EstimateTokens(input + chat.Context.Prompt() + output)),
	}

	_, err = lf.GenerationEnd(g)
	if err != nil {
		chat.Logger.Error().Err(err).Msg("failed to end langfuse generation")
	}

	_, err = lf.Trace(&model.Trace{
		ID:     trace.ID,
		Output: output,
	})
	if err != nil {
		chat.Logger.Error().Err(err).Msg("failed to update langfuse trace")
	}

	return output, nil
}
