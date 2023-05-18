package adapters

import (
	"context"
	"errors"
	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestChatGPTAdapter_CreateImagePrompt(t *testing.T) {
	testCases := []struct {
		name          string
		responseBody  string
		responseCode  int
		expectedError error
	}{
		{
			name: "Success",
			responseBody: `{
				"choices": [
					{
						"message": {
							"role": "assistant",
							"content": "Test Prompt"
						}
					}
				]
			}`,
			responseCode:  http.StatusOK,
			expectedError: nil,
		},
		{
			name:          "API Error",
			responseBody:  "",
			responseCode:  http.StatusInternalServerError,
			expectedError: errors.New("failed to generate image prompt, status code: 500"),
		},
		{
			name:          "Empty Choices",
			responseBody:  `{"choices": []}`,
			responseCode:  http.StatusOK,
			expectedError: errors.New("no choices returned from ChatGPT API"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockClient := newMockHttpClient(t)
			mockClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: tc.responseCode,
				Body:       ioutil.NopCloser(strings.NewReader(tc.responseBody)),
			}, nil)

			adapter := NewChatGPTAdapter("test-api-key", mockClient)

			_, err := adapter.CreateImagePrompt(context.Background(), domain.NewsArticle{Title: "test", Body: "test"})
			if (err != nil && tc.expectedError == nil) ||
				(err == nil && tc.expectedError != nil) ||
				(err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}
