package adapters

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// TestNYTimesAdapter_GetMainArticle tests the GetMainArticle method of the New York Times adapter.
func TestNYTimesAdapter_GetMainArticle(t *testing.T) {
	testCases := []struct {
		name          string
		responseBody  string
		responseCode  int
		expectedError error
	}{
		{
			name: "Success",
			responseBody: `{
				"results": [
					{
						"title": "Test Title",
						"abstract": "Test Abstract",
						"published_date": "2022-01-01T00:00:00-05:00"
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
			expectedError: errors.New("failed to fetch data from New York Times API"),
		},
		{
			name:          "Empty Results",
			responseBody:  `{"results": []}`,
			responseCode:  http.StatusOK,
			expectedError: errors.New("no articles found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock HTTP client.
			mockClient := NewMockClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: tc.responseCode,
					Body:       ioutil.NopCloser(strings.NewReader(tc.responseBody)),
				}
			})

			// Create the adapter with the mock client.
			adapter := &nyTimesAdapter{
				apiKey: "test-api-key",
			}
			adapter.client = mockClient

			// Call GetMainArticle and check for the expected error.
			_, err := adapter.GetMainArticle(context.Background())
			if (err != nil && tc.expectedError == nil) ||
				(err == nil && tc.expectedError != nil) ||
				(err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

// MockClient is a custom HTTP client for testing purposes.
type MockClient struct {
	DoFunc func(req *http.Request) *http.Response
}

// Do is the implementation of the Do method for the custom HTTP client.
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req), nil
}

// NewMockClient creates a new custom HTTP client.
func NewMockClient(fn func(req *http.Request) *http.Response) *MockClient {
	return &MockClient{
		DoFunc: fn,
	}
}
