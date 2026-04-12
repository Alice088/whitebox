package deepseek

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"whitebox/internal/core/llm"
	http2 "whitebox/internal/http"
	"whitebox/internal/providers"
)

type DeepSeek struct {
	baseURL string
	apiKey  string
	model   string
}

func New(opts providers.InitOpts) llm.LLM {
	url := "https://api.deepseek.com"
	if len(opts.BaseURL) != 0 {
		url = opts.BaseURL
	}

	return DeepSeek{
		baseURL: url,
		apiKey:  opts.ApiKey,
		model:   opts.Model,
	}
}

func (d DeepSeek) Ask(prompt string) (string, error) {
	url := d.baseURL + "/v1/chat/completions"

	reqBody := http2.RequestBody{
		Model: d.model,
		Messages: []http2.Message{
			{Role: "system", Content: "You are a helpful assistant."},
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

	client := &http.Client{Timeout: 3 * time.Minute}

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
		return "", fmt.Errorf("status code is %s(%d): %s", resp.Status, resp.StatusCode, string(body))
	}

	var response http2.ResponseBody
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no answer")
}

func (d DeepSeek) Model() string {
	return d.model
}

func (d DeepSeek) EstimateTokens(input string) float64 {
	return float64(len([]rune(input))) * 0.3
}
