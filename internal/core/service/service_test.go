package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/BaronBonet/content-generator/internal/infrastructure"
	"github.com/BaronBonet/go-logger/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
)

func TestService_GenerateNewsContent(t *testing.T) {
	testLogger := logger.NewTestLogger()
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
				newsArticle := domain.NewsArticle{Title: "Test Article", Body: "Test body"}
				mockNewsAdapter.On("GetMainArticle", mock.AnythingOfType("*context.emptyCtx")).Return(newsArticle, nil)
				prompt := fmt.Sprintf("Generate a single sentence image prompt based on the following news title and body:"+
					"\nTitle: %s"+
					"\nBody: %s"+
					"\n Do not include prompts that will be rejected by the Dalle safety system. For example mentioning dictators like Vladimir Putin."+
					"\n\n Examples of good prompts"+
					"\n- 3D render of a pink balloon dog in a violet room"+
					"\n- Illustration of a happy cat sitting on a couch in a living room with a coffee mug in its hand", newsArticle.Title, newsArticle.Body)

				llmAdapter.On("Chat", mock.AnythingOfType("*context.emptyCtx"), prompt).Return(prompt, nil)
				imagePath := "https://test.com/test.jpg"
				generatorName := "TestGenerator"
				mockImageGenerationAdapter.On("GenerateImage", mock.AnythingOfType("*context.emptyCtx"), prompt).Return(domain.ImagePath(imagePath), nil)
				mockImageGenerationAdapter.On("GetGeneratorName").Return(generatorName)
				mockSocialMediaAdapter.On("PublishImagePost", mock.AnythingOfType("*context.emptyCtx"), domain.ImagePath(imagePath), prompt, generatorName, newsArticle).Return(nil)
				mockSocialMediaAdapter.On("GetName").Return("Twitter")
			},
			expectedError: nil,
		},
		{
			name: "NewsAdapterError",
			setupMocks: func() {
				mockNewsAdapter.On("GetMainArticle", mock.Anything).Return(domain.NewsArticle{}, errors.New("news error"))
			},
			expectedError: errors.New("news error"),
		},
		{
			name: "LLMAdapterError",
			setupMocks: func() {
				mockNewsAdapter.On("GetMainArticle", mock.Anything).Return(domain.NewsArticle{Title: "Test Article"}, nil)
				llmAdapter.On("Chat", mock.Anything, mock.Anything).Return("", errors.New("prompt error"))
			},
			expectedError: errors.New("prompt error"),
		},
		{
			name: "ImageGenerationAdapterError",
			setupMocks: func() {
				mockNewsAdapter.On("GetMainArticle", mock.Anything).Return(domain.NewsArticle{Title: "Test Article"}, nil)
				llmAdapter.On("Chat", mock.Anything, mock.Anything).Return("Test Image Prompt", nil)
				mockImageGenerationAdapter.On("GenerateImage", mock.Anything, mock.Anything).Return(domain.ImagePath(""), errors.New("generation error"))
			},
			expectedError: errors.New("generation error"),
		},
		{
			name: "SocialMediaAdapterError",
			setupMocks: func() {
				mockNewsAdapter.On("GetMainArticle", mock.Anything).Return(domain.NewsArticle{Title: "Test Article"}, nil)
				llmAdapter.On("Chat", mock.Anything, mock.Anything).Return("Test Image Prompt", nil)
				mockImageGenerationAdapter.On("GenerateImage", mock.Anything, mock.Anything).Return(domain.ImagePath("Test Image Path"), nil)
				mockImageGenerationAdapter.On("GetGeneratorName").Return("TestGenerator")
				mockSocialMediaAdapter.On("PublishImagePost", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("social media error"))
				mockSocialMediaAdapter.On("GetName").Return("Twitter")
			},
			// We don't want it to retry if the social media adapters fails
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			srv := NewNewsContentService(
				testLogger,
				mockNewsAdapter,
				llmAdapter,
				mockImageGenerationAdapter,
				[]ports.SocialMediaAdapter{mockSocialMediaAdapter},
			)

			err := srv.GenerateNewsContent(context.Background())

			assert.Equal(t, tc.expectedError, err)

			infrastructure.TearDownAdapters(&mockNewsAdapter.Mock,
				&llmAdapter.Mock,
				&mockImageGenerationAdapter.Mock,
				&mockSocialMediaAdapter.Mock,
			)
		})
	}
}
