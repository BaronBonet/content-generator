package adapters

import (
	"bytes"
	"context"
	"fmt"
	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/BaronBonet/content-generator/internal/infrastructure"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestTwitterAdapter_PublishImagePost(t *testing.T) {
	mockOAuthClient := newMockHttpClient(t)
	mockClient := newMockHttpClient(t)
	mockLogger := ports.NewMockLogger(t)

	type testCase struct {
		name          string
		setupMocks    func(*testCase)
		errorResponse error
		newsArticle   domain.NewsArticle
		prompt        string
	}

	testCases := []testCase{
		{
			name: "Download Image Error",
			setupMocks: func(tc *testCase) {
				mockClient.On("Get", mock.AnythingOfType("string")).Return(&http.Response{Body: io.NopCloser(bytes.NewReader([]byte(""))),
					StatusCode: http.StatusInternalServerError,
				}, nil)
			},
			errorResponse: fmt.Errorf("failed to download image, status code: 500"),
		},
		{
			name: "Download Image Error",
			setupMocks: func(tc *testCase) {
				mockClient.On("Get", mock.AnythingOfType("string")).Return(&http.Response{Body: io.NopCloser(bytes.NewReader([]byte(""))),
					StatusCode: http.StatusOK,
				}, nil)
				mockOAuthClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewReader([]byte(tc.errorResponse.Error()))),
				}, nil)
			},
			errorResponse: fmt.Errorf("failed to download image, status code: 500"),
		},
		{
			name: "Tweet Truncated",
			setupMocks: func(tc *testCase) {
				mockClient.On("Get", mock.AnythingOfType("string")).Return(&http.Response{Body: io.NopCloser(bytes.NewReader([]byte(""))),
					StatusCode: http.StatusOK,
				}, nil).Once()
				mockOAuthClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"media_id_string": "12345"}`)),
				}, nil).Once()
				mockLogger.On("Warn", "Tweet was truncated to 280 characters", "full tweet", tc.newsArticle.Title+" "+tc.newsArticle.Url).Once()
				mockOAuthClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"data": {"id": "67890", "text": "test text"}}`))),
				}, nil).Once()
				mockOAuthClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}, nil).Once()
			},
			errorResponse: nil,
			newsArticle: domain.NewsArticle{
				Title: strings.Repeat("a", 300),
				Url:   "https://example.com",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks(&tc)
			twitterAdapter := NewTwitterSocialMediaAdapter(mockOAuthClient, mockClient, mockLogger)
			err := twitterAdapter.PublishImagePost(context.Background(), "https://test.com/test.png", tc.prompt, "test generator", tc.newsArticle)
			require.Equal(t, tc.errorResponse, err)
			mockOAuthClient.AssertExpectations(t)
			mockClient.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			infrastructure.TearDownAdapters(&mockOAuthClient.Mock, &mockClient.Mock, &mockLogger.Mock)
		})
	}
}
