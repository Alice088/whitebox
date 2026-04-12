package pipeline

import (
	"context"
	"errors"
	"time"

	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
	"github.com/rs/zerolog"
)

type State struct {
	IO            IO
	Runtime       Runtime
	Observability Observability
}

type Step func(ctx context.Context, s *State) error

type namedStep struct {
	name string
	fn   Step
}

type Runner struct {
	steps  []namedStep
	Logger *zerolog.Logger
}

func (r *Runner) Use(name string, step Step) {
	r.steps = append(r.steps, namedStep{name: name, fn: step})
}

func (r *Runner) Run(ctx context.Context, state *State) error {
	if state == nil {
		return errors.New("state is nil")
	}

	for _, step := range r.steps {
		if r.Logger != nil {
			r.Logger.Info().Str("step", step.name).Msg("run step")
		}

		if err := step.fn(ctx, state); err != nil {
			return err
		}
	}

	return nil
}

func AskLLM() Step {
	return func(_ context.Context, s *State) error {
		if s.Runtime.LLM == nil {
			return errors.New("llm is nil")
		}

		prompt := s.Runtime.Context.Prompt()
		output, err := s.Runtime.LLM.Ask(prompt)
		if err != nil {
			return err
		}

		s.IO.Output = output
		return nil
	}
}

func Logging(logger zerolog.Logger) Step {
	return func(_ context.Context, s *State) error {
		logger.Info().
			Str("input", s.IO.Input).
			Str("output", s.IO.Output).
			Msg("pipeline log")
		return nil
	}
}

func LangfuseFlush(client *langfuse.Langfuse) Step {
	return func(ctx context.Context, s *State) error {
		client.Flush(ctx)
		return nil
	}
}

func LangfuseStart(client *langfuse.Langfuse) Step {
	return func(_ context.Context, s *State) error {
		if client == nil {
			return nil
		}

		trace, err := client.Trace(&model.Trace{
			Name:      "whitebox-request",
			Input:     s.IO.Input,
			Timestamp: new(time.Now()),
		})
		if err != nil {
			return err
		}

		s.Observability.TraceID = trace.ID
		return nil
	}
}

func LangfuseEnd(client *langfuse.Langfuse) Step {
	return func(_ context.Context, s *State) error {
		if client == nil {
			return nil
		}

		trace, err := client.Trace(&model.Trace{
			ID:     s.Observability.TraceID,
			Output: s.IO.Output,
		})
		if err != nil {
			return err
		}

		s.Observability.TraceID = trace.ID
		return nil
	}
}

func LangfuseGenerationEnd(client *langfuse.Langfuse) Step {
	return func(_ context.Context, s *State) error {
		if client == nil || s.Observability.TraceID == "" {
			return nil
		}

		s.Observability.Generation.Output = model.M{"completion": s.IO.Output}
		s.Observability.Generation.Usage = model.Usage{
			Input:  int(s.Runtime.LLM.EstimateTokens(s.Runtime.Context.Prompt())),
			Output: int(s.Runtime.LLM.EstimateTokens(s.IO.Output)),
			Total:  int(s.Runtime.LLM.EstimateTokens(s.Runtime.Context.Prompt() + s.IO.Output)),
		}
		_, err := client.GenerationEnd(s.Observability.Generation)
		if err != nil {
			return err
		}

		return nil
	}
}

func LangfuseGenerationStart(client *langfuse.Langfuse) Step {
	return func(_ context.Context, s *State) error {
		if client == nil || s.Observability.TraceID == "" {
			return nil
		}

		g, err := client.Generation(&model.Generation{
			Model:   s.Runtime.LLM.Model(),
			Name:    "llm-call",
			TraceID: s.Observability.TraceID,
			Input: []model.M{
				{"role": "user", "content": s.Runtime.Context.Prompt()},
			},
		}, nil)
		if err != nil {
			return err
		}

		s.Observability.Generation = g

		return nil
	}
}
