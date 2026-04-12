package llamacpp

import (
	"bytes"
	"coreclaw/internal/context"
	"coreclaw/internal/llm"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
	"github.com/rs/zerolog"
)

type Model struct {
	baseURL  string
	apiKey   string
	model    string
	langFuse *langfuse.Langfuse
	logger   *zerolog.Logger
	context  context.Context
}

type Opts struct {
	ApiKey   string
	BaseURL  string
	Model    string
	LangFuse *langfuse.Langfuse
	Logger   *zerolog.Logger
	Context  context.Context
}

func New(opts Opts) llm.LLM {
	url := "http://0.0.0.0:8080"
	if len(opts.BaseURL) != 0 {
		url = opts.BaseURL
	}

	return Model{
		baseURL:  url,
		apiKey:   opts.ApiKey,
		model:    opts.Model,
		langFuse: opts.LangFuse,
		logger:   opts.Logger,
		context:  opts.Context,
	}
}

func (d Model) Ask(prompt string, id string) (string, error) {
	url := d.baseURL + "/v1/chat/completions"

	reqBody := llm.RequestBody{
		Model: d.model,
		Messages: []llm.Message{
			{Role: "system", Content: d.context.Prompt()},
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+d.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 3 * time.Minute,
	}

	g, err := d.langFuse.Generation(&model.Generation{
		Model:   d.model,
		Name:    "llm-call",
		TraceID: id,
		Input: []model.M{
			{"role": "system", "content": d.context.Prompt()},
			{"role": "user", "content": prompt},
		},
	}, nil)
	if err != nil {
		return "", err
	}

	var answer string
	defer func() {
		g.Output = model.M{"completion": answer}
		g.Usage = model.Usage{
			Input:  int(d.EstimateTokens(prompt)),
			Output: int(d.EstimateTokens(answer)),
			Total:  int(d.EstimateTokens(answer + prompt + d.context.Prompt())),
		}
		_, gErr := d.langFuse.GenerationEnd(g)
		if gErr != nil {
			d.logger.Error().Err(gErr).Msg("Failed to generation_end")
		}
	}()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", string(body))
		return "", fmt.Errorf("status code is %s(%d): %s", resp.Status, resp.StatusCode, string(body))
	}

	var response llm.ResponseBody
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	if len(response.Choices) > 0 {
		answer = response.Choices[0].Message.Content
		return answer, nil
	}

	return "", fmt.Errorf("no answer")
}

func (d Model) Model() string {
	return d.model
}

func (d Model) EstimateTokens(input string) float64 {
	return float64(len([]rune(input))) * 1
}
