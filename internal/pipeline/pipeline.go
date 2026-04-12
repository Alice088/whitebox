package pipeline

import (
	"context"
	"errors"
	"time"
	syscontext "whitebox/internal/context"
	llmcore "whitebox/internal/core/llm"

	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
	"github.com/rs/zerolog"
)

type State struct {
	Name         string
	Input        string
	SystemPrompt string
	Output       string

	LLM     llmcore.LLM
	Context syscontext.Context

	TraceID string
	Meta    map[string]any
}

type ReadStep func(ctx context.Context, s State) error
type MutStep func(ctx context.Context, s *State) error

type Runner struct {
	R      map[string]ReadStep
	W      map[string]MutStep
	Logger *zerolog.Logger
}

func (r *Runner) Read(step ReadStep, name string) {
	r.R[name] = step
}

func (r *Runner) Write(step MutStep, name string) {
	r.W[name] = step
}

func (r *Runner) Run(ctx context.Context, state *State) error {
	if state == nil {
		return errors.New("state is nil")
	}

	for name, step := range r.R {
		r.Logger.Info().Str("step_mode", "read").Str("step_name", name).Msg("run step")

		if err := step(ctx, *state); err != nil {
			return err
		}
	}

	for name, step := range r.W {
		r.Logger.Info().Str("step_mode", "write").Str("step_name", name).Msg("run step")

		if err := step(ctx, state); err != nil {
			return err
		}
	}

	return nil
}

func BuildPrompt() MutStep {
	return func(_ context.Context, s *State) error {
		s.SystemPrompt = s.Context.Prompt()
		return nil
	}
}

func AskLLM() MutStep {
	return func(_ context.Context, s *State) error {
		if s.LLM == nil {
			return errors.New("llm is nil")
		}

		output, err := s.LLM.Ask(s.Input, s.SystemPrompt)
		if err != nil {
			return err
		}

		s.Output = output
		return nil
	}
}

func Logging(logger zerolog.Logger) ReadStep {
	return func(_ context.Context, s State) error {
		logger.Info().
			Str("input", s.Input).
			Str("output", s.Output).
			Msg("pipeline log")
		return nil
	}
}

func LangfuseStart(client *langfuse.Langfuse, name string) MutStep {
	return func(_ context.Context, s *State) error {
		if client == nil {
			return nil
		}

		trace, err := client.Trace(&model.Trace{
			Name:      name,
			Input:     s.Input,
			Timestamp: new(time.Now()),
		})
		if err != nil {
			return err
		}

		s.TraceID = trace.ID
		return nil
	}
}

func LangfuseEnd(client *langfuse.Langfuse) MutStep {
	return func(_ context.Context, s *State) error {
		if client == nil || s.TraceID == "" {
			return nil
		}

		_, err := client.Trace(&model.Trace{
			ID:     s.TraceID,
			Output: s.Output,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
