package adapters

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestDalleAdapter_GenerateImage(t *testing.T) {
	testCases := []struct {
		name             string
		mockResponse     string
		mockResponseCode int
		expectedError    string
	}{
		{
			name:             "Success",
			mockResponse:     `{"id":[{"url":"http://example.com/image1.png"}]}`,
			mockResponseCode: http.StatusOK,
			expectedError:    "",
		},
		{
			name:             "API Error",
			mockResponse:     `{"error":"API Error"}`,
			mockResponseCode: http.StatusInternalServerError,
			expectedError:    "failed to generate image, status code: 500",
		},
		{
			name:             "Empty Choices",
			mockResponse:     `{"id":[]}`,
			mockResponseCode: http.StatusOK,
			expectedError:    "no choices returned from ChatGPT API",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockClient := newMockHttpClient(t)
			mockClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: tc.mockResponseCode,
				Body:       ioutil.NopCloser(strings.NewReader(tc.mockResponse)),
			}, nil)

			dalleAdapter := NewDalleAdapter("test-api-key", mockClient)

			_, err := dalleAdapter.GenerateImage(context.Background(), "test-prompt")

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
