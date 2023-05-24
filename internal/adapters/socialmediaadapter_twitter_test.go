package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestTwitterAdapter_PublishImagePost(t *testing.T) {
	testCases := []struct {
		name            string
		prompt          string
		imageUploadResp string
		tweetResp       string
		imageUploadCode int
		tweetCode       int
		expectedError   error
		httpDoCalls     int
	}{
		{
			name:            "Success",
			prompt:          "example prompt",
			imageUploadResp: `{"media_id_string": "12345"}`,
			tweetResp:       "{}",
			imageUploadCode: http.StatusOK,
			tweetCode:       http.StatusCreated,
			expectedError:   nil,
			httpDoCalls:     2,
		},
		{
			name:            "Image Upload Error",
			prompt:          "example prompt",
			imageUploadResp: "failed to upload image",
			tweetResp:       "",
			imageUploadCode: http.StatusInternalServerError,
			tweetCode:       http.StatusOK,
			expectedError:   errors.New("failed to upload image"),
			httpDoCalls:     1,
		},
		{
			name:            "Tweet Post Error",
			prompt:          "example prompt",
			imageUploadResp: `{"media_id_string": "12345"}`,
			tweetResp:       "",
			imageUploadCode: http.StatusOK,
			tweetCode:       http.StatusInternalServerError,
			expectedError:   errors.New("failed to post tweet, status code: 500"),
			httpDoCalls:     2,
		},
		{
			name:            "Prompt Greater Than 280 Characters",
			prompt:          strings.Repeat("A", 300), // 300 > 280
			imageUploadResp: `{"media_id_string": "12345"}`,
			tweetResp:       "{}",
			imageUploadCode: http.StatusOK,
			tweetCode:       http.StatusCreated,
			expectedError:   nil,
			httpDoCalls:     2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockOAuthClient := newMockHttpClient(t)
			mockClient := newMockHttpClient(t)
			mockLogger := ports.NewMockLogger(t)
			mockClient.On("Get", "https://test.com/test.png").Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("test image")),
			}, nil).Once()

			mockOAuthClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: tc.imageUploadCode,
				Body:       ioutil.NopCloser(strings.NewReader(tc.imageUploadResp)),
			}, nil).Once()

			if tc.imageUploadCode == http.StatusOK {
				mockOAuthClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
					reqBytes, _ := ioutil.ReadAll(req.Body)
					req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes)) // Reconstruct req.Body as it has been read
					t := tweet{}
					err := json.Unmarshal(reqBytes, &t)
					if err != nil {
						return false
					}
					if len(t.Text) > 280 {
						mockLogger.On("Warn", "Tweet was truncated to 280 characters", "full tweet", t.Text).Once()
					}
					// Verify the length of the tweet text.
					return t.Text == tc.prompt
				})).Return(&http.Response{
					StatusCode: tc.tweetCode,
					Body:       ioutil.NopCloser(strings.NewReader(tc.tweetResp)),
				}, nil).Once()

				mockOAuthClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
					reqBytes, _ := ioutil.ReadAll(req.Body)
					req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes)) // Reconstruct req.Body as it has been read
					t := replyTweet{}
					err := json.Unmarshal(reqBytes, &t)
					if err != nil {
						return false
					}
					if len(t.Text) > 280 {
						mockLogger.On("Warn", "Tweet was truncated to 280 characters", "full tweet", t.Text).Once()
					}
					// Verify the length of the tweet text.
					return t.Text == fmt.Sprintf("Created by %s with the prompt:\n\n%s", "https://example.com", tc.prompt)
				})).Return(&http.Response{
					StatusCode: tc.tweetCode,
					Body:       ioutil.NopCloser(strings.NewReader(tc.tweetResp)),
				}, nil).Once()
			}

			adapter := NewTwitterSocialMediaAdapter(mockOAuthClient, mockClient, mockLogger)

			err := adapter.PublishImagePost(context.Background(), "https://test.com/test.png", tc.prompt, "dalle","https://example.com")

			if (err != nil && tc.expectedError == nil) ||
				(err == nil && tc.expectedError != nil) ||
				(err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}
			mockOAuthClient.AssertNumberOfCalls(t, "Do", tc.httpDoCalls)

			mockLogger.AssertExpectations(t)
		})
	}
	``
