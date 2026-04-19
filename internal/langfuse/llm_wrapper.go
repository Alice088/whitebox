package langfuse

import (
	"context"
	"errors"
	"fmt"
	"time"
	"whitebox/internal/core/llm"
	"whitebox/internal/core/tools"
	"whitebox/pkg/maps"
	"whitebox/pkg/meta"

	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
)

type LLMWrapper struct {
	LLM      llm.LLM
	Langfuse *langfuse.Langfuse

	first bool
	State map[string]any
}

func (w *LLMWrapper) EndTrace() {
	w.State = nil
}

func (w *LLMWrapper) StartTrace(originInput string) error {
	if w.State == nil {
		w.State = make(map[string]any)
	}

	if w.Langfuse == nil {
		w.Langfuse = langfuse.New(context.Background())
	}

	if maps.Exists(w.State, "trace") {
		return nil
	}

	trace, err := w.Langfuse.Trace(&model.Trace{
		Name:      fmt.Sprintf("whitebox-%s-request", meta.AgentName),
		Input:     originInput,
		Timestamp: new(time.Now()),
	})
	if err != nil {
		return err
	}

	w.first = true
	w.State["trace"] = trace
	return err
}

func (w *LLMWrapper) Ask(prompt, systemPrompt string) (out string, err error) {
	t := w.State["trace"].(*model.Trace)
	t.Metadata = map[string]string{
		"user_message_with_tool_prompt": prompt,
	}
	_, err = w.Langfuse.Trace(t)
	if err != nil {
		return "", fmt.Errorf("failed update trace: %w", err)
	}
	defer func() {
		t.Output = out
		_, err = w.Langfuse.Trace(t)
		if err != nil {
			t.Output = model.M{
				"completion": out,
				"error":      err,
			}
			err = fmt.Errorf("failed update trace: %w", err)
		}
	}()

	start := new(time.Now())
	generation, err := w.Langfuse.Generation(&model.Generation{
		Model:     w.LLM.Model(),
		Name:      "llm-call",
		TraceID:   w.State["trace"].(*model.Trace).ID,
		StartTime: start,
		Input: []model.M{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": prompt},
		},
	}, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create generation: %w", err)
	}

	w.State["generation"] = generation

	defer func() {
		generation.Output = model.M{"completion": out}
		generation.EndTime = new(time.Now())
		generation.Usage = model.Usage{
			Input:  int(w.LLM.EstimateTokens(prompt)),
			Output: int(w.LLM.EstimateTokens(out)),
			Total:  int(w.LLM.EstimateTokens(out + prompt + systemPrompt)),
		}
		_, err = w.Langfuse.GenerationEnd(generation)
		if err != nil {
			generation.Output = model.M{
				"completion": out,
				"error":      err,
			}
			err = fmt.Errorf("failed to end generation: %w", err)
		}
	}()

	out, err = w.LLM.Ask(prompt, systemPrompt)
	if err != nil {
		return "", err
	}

	return out, err
}

func (w *LLMWrapper) Tool(call tools.ToolCall) (out string, err error) {
	if !maps.Exists(w.State, "generation") {
		return "", errors.New("generation not exists")
	}

	start := new(time.Now())
	generation, err := w.Langfuse.Generation(&model.Generation{
		TraceID:   w.State["trace"].(*model.Trace).ID,
		Name:      "tool:" + call.Tool,
		StartTime: start,
		Input: map[string]any{
			"arguments": call.Arguments,
		},
		Metadata: map[string]any{
			"type": "tool_call",
		},
	}, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create span(tool_call): %w", err)
	}

	defer func() {
		generation.Output = model.M{"completion": out}
		generation.EndTime = new(time.Now())
		_, err = w.Langfuse.GenerationEnd(generation)
		if err != nil {
			generation.Output = model.M{
				"completion": out,
				"error":      err,
			}
			err = fmt.Errorf("failed to end span: %w", err)
		}
	}()

	return tools.Execute(call)
}

func (w *LLMWrapper) EstimateTokens(s string) float64 {
	return w.LLM.EstimateTokens(s)
}

func (w *LLMWrapper) Model() string {
	return w.LLM.Model()
}
