package chat

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	syscontext "whitebox/internal/core/context"
	"whitebox/internal/core/llm"
	"whitebox/internal/core/status"
	"whitebox/pkg/messages"

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
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Whitebox Chat Mode")
	fmt.Printf("Session ID: %s\n", chat.Session.ID)
	fmt.Println("Type '@exit' to quit, '@clear' to clear history")
	fmt.Println()

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.Trim(scanner.Text(), "\r\n\t ")

		if strings.Contains(input, "@exit") {
			fmt.Println("Exiting...")
			break
		}
		if strings.Contains(input, "@clear") {
			chat.Context.ClearMessages()
			if err := chat.Session.SaveSession(chat.Context.Messages); err != nil {
				messages.PrintError(fmt.Errorf("failed to save cleared session: %w", err))
				chat.Logger.Error().Err(err).Msg("failed to save cleared session")
			} else {
				chat.Logger.Info().Msg("History cleared and saved.")
				fmt.Println("History cleared and saved.")
			}
			continue
		}

		chat.Context.AddMessage(syscontext.Message{
			Role:    "user",
			Content: input,
		})

		chat.Context.TrimMessages(chat.Session.MaxMessages)

		animation := status.NewAnimationController(chat.statusEngine)
		animation.Start()

		answer, err := chat.ask(ctx, input)

		animation.Stop()

		if err != nil {
			messages.PrintError(err)
			chat.Logger.Error().Err(err).Msg("LLM request failed")
			continue
		}

		messages.PrintAssistant(answer)

		chat.Context.AddMessage(syscontext.Message{
			Role:    "assistant",
			Content: answer,
		})

		chat.Context.TrimMessages(chat.Session.MaxMessages)

		if len(chat.Context.Messages) > 0 {
			if err = chat.Session.SaveSession(chat.Context.Messages); err != nil {
				messages.PrintError(fmt.Errorf("failed to save session: %w", err))
				chat.Logger.Error().Err(err).Msg("failed to save session")
			} else {
				chat.Logger.Info().
					Str("session_path", chat.Session.Path).
					Int("saved_messages", len(chat.Context.Messages)).
					Msg("session saved")
			}
		}
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
