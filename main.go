package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	cfg "whitebox/internal/config"
	syscontext "whitebox/internal/core/context"
	xllm "whitebox/internal/core/llm"
	"whitebox/internal/core/status"
	"whitebox/internal/factory"
	"whitebox/internal/flag"
	"whitebox/internal/providers"

	"github.com/caarlos0/env/v11"
	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	logRotator := &lumberjack.Logger{
		Filename:   "./logs/whitebox.log",
		MaxSize:    10,
		MaxBackups: 2,
		MaxAge:     28,
		Compress:   true,
	}

	logger := zerolog.New(logRotator).With().Timestamp().Logger()

	input, err := flag.ParseFlags()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	err = godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	var config cfg.Config
	err = env.Parse(&config)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	var sessionID string
	if input.SessionID != "" {
		sessionID = input.SessionID
	} else {
		sessionID = syscontext.NewSessionID()
	}

	sessionsDir := "context/sessions"
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		logger.Fatal().Err(err).Msg("failed to create sessions directory")
	}

	sessionPath := filepath.Join(sessionsDir, sessionID+".json")
	msgs, err := syscontext.LoadSession(sessionPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load session")
	}

	logger.Info().
		Str("session_id", sessionID).
		Int("loaded_messages", len(msgs)).
		Msg("session loaded")

	systemContext, err := syscontext.NewDefault()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load context")
	}

	systemContext.Messages = msgs

	llm, err := factory.LLM(input.Provider, providers.InitOpts{
		Model:  input.Model,
		ApiKey: config.LLM.ApiKey,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init LLM")
	}

	statusGen := status.NewStatusGenerator()
	runChat(context.Background(), llm, &systemContext, sessionPath, input.MaxHistory, statusGen, logger)
}

func runChat(ctx context.Context, llm xllm.LLM, systemContext *syscontext.Context,
	sessionPath string, maxHistory int, statusGen *status.StatusGenerator,
	logger zerolog.Logger) {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Whitebox Chat Mode")
	fmt.Printf("Session ID: %s\n", filepath.Base(sessionPath))
	fmt.Println("Type '/exit' to quit, '/clear' to clear history")
	fmt.Println()

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.Trim(scanner.Text(), "\r\n\t ")

		if input == "/exit" {
			fmt.Println("Exiting...")
			break
		}
		if input == "/clear" {
			systemContext.ClearMessages()
			if err := syscontext.SaveSession(sessionPath, systemContext.Messages); err != nil {
				fmt.Printf("\033[91mWarning: failed to save cleared session: %v\033[0m\n", err)
				logger.Error().Err(err).Msg("failed to save cleared session")
			} else {
				fmt.Println("History cleared and saved.")
			}
			continue
		}

		systemContext.AddMessage(syscontext.Message{
			Role:    "user",
			Content: input,
		})

		systemContext.TrimMessages(maxHistory)

		animation := status.NewAnimationController(statusGen)
		animation.Start()

		output, err := askWithLangfuse(ctx, llm, systemContext, input, logger)

		animation.Stop()

		if err != nil {
			fmt.Printf("\033[91mError: %v\033[0m\n", err)
			logger.Error().Err(err).Msg("LLM request failed")
			continue
		}

		fmt.Printf("\x1b[47m  \x1b[0m whitebox >  %s\n", output)

		systemContext.AddMessage(syscontext.Message{
			Role:    "assistant",
			Content: output,
		})

		systemContext.TrimMessages(maxHistory)

		if len(systemContext.Messages) > 0 {
			if err := syscontext.SaveSession(sessionPath, systemContext.Messages); err != nil {
				fmt.Printf("\033[91mWarning: failed to save session: %v\033[0m\n", err)
				logger.Error().Err(err).Msg("failed to save session")
			} else {
				logger.Info().
					Str("session_path", sessionPath).
					Int("saved_messages", len(systemContext.Messages)).
					Msg("session saved")
			}
		}
	}
}

func askWithLangfuse(ctx context.Context, llm xllm.LLM, systemContext *syscontext.Context, input string, logger zerolog.Logger) (string, error) {
	publicKey := os.Getenv("LANGFUSE_PUBLIC_KEY")
	secretKey := os.Getenv("LANGFUSE_SECRET_KEY")
	if publicKey == "" || secretKey == "" {
		return llm.Ask(systemContext.Prompt())
	}

	lf := langfuse.New(ctx)
	defer lf.Flush(ctx)

	trace, err := lf.Trace(&model.Trace{
		Name:      "whitebox-chat",
		Input:     input,
		Timestamp: new(time.Now()),
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create langfuse trace")
		return llm.Ask(systemContext.Prompt())
	}

	g, err := lf.Generation(&model.Generation{
		Model:   llm.Model(),
		Name:    "llm-call",
		TraceID: trace.ID,
		Input: []model.M{
			{"role": "user", "content": systemContext.Prompt()},
		},
	}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create langfuse generation")
		return llm.Ask(systemContext.Prompt())
	}

	output, err := llm.Ask(systemContext.Prompt())
	if err != nil {
		_, gErr := lf.GenerationEnd(&model.Generation{
			ID:     g.ID,
			Output: model.M{"error": err.Error()},
		})
		if gErr != nil {
			logger.Error().Err(err).Msg("failed to end langfuse generation")
		}
		return "", err
	}

	g.Output = model.M{"completion": output}
	g.Usage = model.Usage{
		Input:  int(llm.EstimateTokens(systemContext.Prompt())),
		Output: int(llm.EstimateTokens(output)),
		Total:  int(llm.EstimateTokens(systemContext.Prompt() + output)),
	}

	_, err = lf.GenerationEnd(g)
	if err != nil {
		logger.Error().Err(err).Msg("failed to end langfuse generation")
	}

	_, err = lf.Trace(&model.Trace{
		ID:     trace.ID,
		Output: output,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to update langfuse trace")
	}

	return output, nil
}
