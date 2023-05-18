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
	// CreateImagePrompt creates a prompt that an AI image generator can use
	CreateImagePrompt(ctx context.Context, article domain.NewsArticle) (domain.ImagePrompt, error)
}

// ImageGenerationAdapter is responsible for connecting to image generation models like DALL-E, Midjourney or Stable Diffusion
//
//go:generate mockery --name=ImageGenerationAdapter
type ImageGenerationAdapter interface {
	// GenerateImage generates an image from a prompt
	GenerateImage(ctx context.Context, prompt domain.ImagePrompt) (domain.ImagePath, error)
}

// SocialMediaAdapter is responsible for connecting to social media services like Twitter
//
//go:generate mockery --name=SocialMediaAdapter
type SocialMediaAdapter interface {
	// PublishImagePost publishes an image post to a social media service
	PublishImagePost(ctx context.Context, image domain.ImagePath, prompt domain.ImagePrompt, sourceUrl string) error
}
