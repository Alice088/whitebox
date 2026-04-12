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
	Input        string
	SystemPrompt string
	Output       string

	LLM     llmcore.LLM
	Context syscontext.Context

	TraceID string
	Meta    map[string]any
}

type Step func(ctx context.Context, s *State) error

type Runner struct {
	steps []Step
}

func (r *Runner) Use(step Step) {
	r.steps = append(r.steps, step)
}

func (r *Runner) Run(ctx context.Context, state *State) error {
	if state == nil {
		return errors.New("state is nil")
	}

	for _, step := range r.steps {
		if err := step(ctx, state); err != nil {
			return err
		}
	}

	return nil
}

func BuildPrompt() Step {
	return func(_ context.Context, s *State) error {
		s.SystemPrompt = s.Context.Prompt()
		return nil
	}
}

func AskLLM() Step {
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

func Logging(logger zerolog.Logger) Step {
	return func(_ context.Context, s *State) error {
		logger.Info().
			Str("input", s.Input).
			Str("output", s.Output).
			Msg("pipeline log")
		return nil
	}
}

func LangfuseStart(client *langfuse.Langfuse, name string) Step {
	return func(_ context.Context, s *State) error {
		if client == nil {
			return nil
		}

		now := time.Now()

		trace, err := client.Trace(&model.Trace{
			Name:      name,
			Input:     s.Input,
			Timestamp: &now,
		})
		if err != nil {
			return err
		}

		s.TraceID = trace.ID
		return nil
	}
}

func LangfuseEnd(client *langfuse.Langfuse) Step {
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
