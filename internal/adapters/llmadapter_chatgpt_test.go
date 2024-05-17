package adapters

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
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
				Body:       io.NopCloser(strings.NewReader(tc.responseBody)),
			}, nil)

			adapter := NewChatGPTAdapter("test-api-key", mockClient)

			_, err := adapter.Chat(context.Background(), "Test prompt")
			if (err != nil && tc.expectedError == nil) ||
				(err == nil && tc.expectedError != nil) ||
				(err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}
