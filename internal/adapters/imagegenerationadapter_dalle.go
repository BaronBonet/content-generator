package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"io/ioutil"
	"net/http"
)

const dalleAPIURL = "https://api.openai.com/v1/images/generations"

type dalleAdapter struct {
	apiKey string
	client httpClient
}

func NewDalleAdapter(apiKey string, httpClient httpClient) ports.ImageGenerationAdapter {
	return &dalleAdapter{
		apiKey: apiKey,
		client: httpClient,
	}
}

func (d *dalleAdapter) GenerateImage(ctx context.Context, prompt domain.ImagePrompt) (domain.ImagePath, error) {
	requestBody := map[string]interface{}{
		"prompt": prompt,
		"n":      1,
		"size":   "256x256",
	}
	jsonRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, dalleAPIURL, bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.apiKey))

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to generate image, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	response := dalleApiResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body: %w", err)
	}

	if len(response.Choices) > 0 {
		return domain.ImagePath(response.Choices[0].Url), nil
	} else {
		return "", fmt.Errorf("no choices returned from ChatGPT API")
	}
}

type dalleData struct {
	Url string `json:"url"`
}

type dalleApiResponse struct {
	Choices []dalleData `json:"id"`
}
