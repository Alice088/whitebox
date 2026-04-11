package deepseek

import (
	"bytes"
	"coreclaw/internal/llm"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DeepSeek struct {
	baseURL string
	apiKey  string
	model   string
}

type Opts struct {
	ApiKey  string
	BaseURL string
	Model   string
}

func New(opts Opts) llm.LLM {
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

func (d DeepSeek) Ask(s string) (string, error) {
	url := d.baseURL + "/v1/chat/completions"

	reqBody := llm.RequestBody{
		Model: d.model,
		Messages: []llm.Message{
			{Role: "system", Content: "You are a helpful coding assistant"},
			{Role: "user", Content: s},
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
		Timeout: 10 * time.Second,
	}
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
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no answer")
}

func (d DeepSeek) Model() string {
	return d.model
}
