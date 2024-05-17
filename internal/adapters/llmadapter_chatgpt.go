package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/BaronBonet/content-generator/internal/core/ports"
)

const chatGPTAPIURL = "https://api.openai.com/v1/chat/completions"

type chatGPTAdapter struct {
	client httpClient
	apiKey string
}

func NewChatGPTAdapter(apiKey string, httpClient httpClient) ports.LLMAdapter {
	return &chatGPTAdapter{
		apiKey: apiKey,
		client: httpClient,
	}
}

func (c *chatGPTAdapter) Chat(ctx context.Context, prompt string) (string, error) {

	requestBody := map[string]interface{}{
		"model":       "gpt-4o",
		"temperature": 0.9,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}
	jsonRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, chatGPTAPIURL, bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to generate image prompt, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body: %w", err)
	}

	if len(apiResponse.Choices) > 0 {
		return apiResponse.Choices[0].Message.Content, nil
	} else {
		return "", fmt.Errorf("no choices returned from ChatGPT API")
	}
}
