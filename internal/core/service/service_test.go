package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
)

// tearDownAdapters resets the mocks so that they can be reused
func tearDownAdapters(adapters ...*mock.Mock) {
	for _, adapter := range adapters {
		adapter.ExpectedCalls = []*mock.Call{}
		adapter.Calls = []mock.Call{}
	}
}

func TestService_GenerateNewsContent(t *testing.T) {
	mockLogger := ports.NewMockLogger(t)
	mockNewsAdapter := ports.NewMockNewsAdapter(t)
	llmAdapter := ports.NewMockLLMAdapter(t)
	mockImageGenerationAdapter := ports.NewMockImageGenerationAdapter(t)
	mockSocialMediaAdapter := ports.NewMockSocialMediaAdapter(t)

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				mockNewsAdapter.On("GetMainArticle", mock.Anything).Return(domain.NewsArticle{Title: "Test Article"}, nil)
				llmAdapter.On("CreateImagePrompt", mock.Anything, mock.Anything).Return(domain.ImagePrompt("Test Image Prompt"), nil)
				mockImageGenerationAdapter.On("GenerateImage", mock.Anything, mock.Anything).Return(domain.ImagePath("Test Image Path"), nil)
				mockSocialMediaAdapter.On("PublishImagePost", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "NewsAdapterError",
			setupMocks: func() {
				mockNewsAdapter.On("GetMainArticle", mock.Anything).Return(domain.NewsArticle{}, errors.New("news error"))
				mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: errors.New("news error"),
		},
		{
			name: "LLMAdapterError",
			setupMocks: func() {
				mockNewsAdapter.On("GetMainArticle", mock.Anything).Return(domain.NewsArticle{Title: "Test Article"}, nil)
				llmAdapter.On("CreateImagePrompt", mock.Anything, mock.Anything).Return(domain.ImagePrompt(""), errors.New("prompt error"))
				mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: errors.New("prompt error"),
		},
		{
			name: "ImageGenerationAdapterError",
			setupMocks: func() {
				mockNewsAdapter.On("GetMainArticle", mock.Anything).Return(domain.NewsArticle{Title: "Test Article"}, nil)
				llmAdapter.On("CreateImagePrompt", mock.Anything, mock.Anything).Return(domain.ImagePrompt("Test Image Prompt"), nil)
				mockImageGenerationAdapter.On("GenerateImage", mock.Anything, mock.Anything).Return(domain.ImagePath(""), errors.New("generation error"))
				mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: errors.New("generation error"),
		},
		{
			name: "SocialMediaAdapterError",
			setupMocks: func() {
				mockNewsAdapter.On("GetMainArticle", mock.Anything).Return(domain.NewsArticle{Title: "Test Article"}, nil)
				llmAdapter.On("CreateImagePrompt", mock.Anything, mock.Anything).Return(domain.ImagePrompt("Test Image Prompt"), nil)
				mockImageGenerationAdapter.On("GenerateImage", mock.Anything, mock.Anything).Return(domain.ImagePath("Test Image Path"), nil)
				mockSocialMediaAdapter.On("PublishImagePost", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("social media error"))
				mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: errors.New("social media error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			srv := NewNewsContentService(
				mockLogger,
				mockNewsAdapter,
				llmAdapter,
				mockImageGenerationAdapter,
				mockSocialMediaAdapter,
			)

			err := srv.GenerateNewsContent(context.Background())

			assert.Equal(t, tc.expectedError, err)

			tearDownAdapters(&mockNewsAdapter.Mock,
				&llmAdapter.Mock,
				&mockImageGenerationAdapter.Mock,
				&mockSocialMediaAdapter.Mock,
				&mockLogger.Mock)
		})
	}
}
