package adapters

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
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
						"published_date": "2022-01-01T00:00:00-05:00",
						"url": "https://www.nytimes.com/2022/01/01/test-article.html"
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
			mockClient := newMockHttpClient(t)

			mockClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: tc.responseCode,
				Body:       ioutil.NopCloser(strings.NewReader(tc.responseBody)),
			}, nil)

			adapter := NewNYTimesNewsAdapter("test-api-key", mockClient)

			_, err := adapter.GetMainArticle(context.Background())
			if (err != nil && tc.expectedError == nil) ||
				(err == nil && tc.expectedError != nil) ||
				(err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}

			mockClient.AssertCalled(t, "Do", mock.Anything)
		})
	}
}
