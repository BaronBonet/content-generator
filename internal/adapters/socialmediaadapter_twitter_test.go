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

// TestTwitterAdapter_PublishImagePost tests the PublishImagePost method of the Twitter adapter.
func TestTwitterAdapter_PublishImagePost(t *testing.T) {
	testCases := []struct {
		name            string
		imageUploadResp string
		tweetResp       string
		imageUploadCode int
		tweetCode       int
		expectedError   error
		httpDoCalls     int
	}{
		{
			name:            "Success",
			imageUploadResp: `{"media_id_string": "12345"}`,
			tweetResp:       "{}",
			imageUploadCode: http.StatusOK,
			tweetCode:       http.StatusCreated,
			expectedError:   nil,
			httpDoCalls:     2,
		},
		{
			name:            "Image Upload Error",
			imageUploadResp: "failed to upload image",
			tweetResp:       "",
			imageUploadCode: http.StatusInternalServerError,
			tweetCode:       http.StatusOK,
			expectedError:   errors.New("failed to upload image"),
			httpDoCalls:     1,
		},
		{
			name:            "Tweet Post Error",
			imageUploadResp: `{"media_id_string": "12345"}`,
			tweetResp:       "",
			imageUploadCode: http.StatusOK,
			tweetCode:       http.StatusInternalServerError,
			expectedError:   errors.New("failed to post tweet, status code: 500"),
			httpDoCalls:     2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := newMockHttpClient(t)

			mockClient.On("Get", "https://test.com/test.png").Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("test image")),
			}, nil).Once()

			mockClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: tc.imageUploadCode,
				Body:       ioutil.NopCloser(strings.NewReader(tc.imageUploadResp)),
			}, nil).Once()

			if tc.imageUploadCode == http.StatusOK {
				mockClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: tc.tweetCode,
					Body:       ioutil.NopCloser(strings.NewReader(tc.tweetResp)),
				}, nil).Once()
			}

			adapter := NewTwitterAdapter(mockClient)

			err := adapter.PublishImagePost(context.Background(), "https://test.com/test.png", "example prompt", "https://example.com")
			if (err != nil && tc.expectedError == nil) ||
				(err == nil && tc.expectedError != nil) ||
				(err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
