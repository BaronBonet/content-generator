package ports

import (
	"context"
	"github.com/BaronBonet/content-generator/internal/core/domain"
)

// NewsAdapter interacts with external news services
//
//go:generate mockery --name=NewsAdapter
type NewsAdapter interface {
	// GetMainArticle finds the main article, the concept of the main article will be adapter specific.
	GetMainArticle(ctx context.Context) (domain.NewsArticle, error)
}

// LLMAdapter is responsible for connecting to large language models like ChatGPT
//
//go:generate mockery --name=LLMAdapter
type LLMAdapter interface {
	// Chat has a conversation with a large language model
	Chat(ctx context.Context, prompt string) (string, error)
}

// ImageGenerationAdapter is responsible for connecting to image generation models like DALL-E, Midjourney or Stable Diffusion
//
//go:generate mockery --name=ImageGenerationAdapter
type ImageGenerationAdapter interface {
	// GenerateImage generates an image from a prompt
	GenerateImage(ctx context.Context, prompt string) (domain.ImagePath, error)
	// GetGeneratorName returns the name of the generator e.g. "DALL-E" or "Midjourney"
	GetGeneratorName() string
}

// SocialMediaAdapter is responsible for connecting to social media services like Twitter
//
//go:generate mockery --name=SocialMediaAdapter
type SocialMediaAdapter interface {
	// PublishImagePost publishes an image post to a social media service
	PublishImagePost(ctx context.Context, image domain.ImagePath, prompt string, imageGeneratorName string, newsArticle domain.NewsArticle) error
}
