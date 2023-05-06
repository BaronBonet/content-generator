package ports

import (
	"context"
	"github.com/BaronBonet/content-generator/internal/core/domain"
)

// NewsAdapter interacts with external news services
type NewsAdapter interface {
	// GetMainArticle finds the main article, the concept of the main article will be adapter specific.
	GetMainArticle(ctx context.Context) (domain.NewsArticle, error)
}

type PromptCreationAdapter interface {
	CreateImagePrompt(ctx context.Context, article domain.NewsArticle) (domain.ImagePrompt, error)
}

type ImageGenerationAdapter interface {
	GenerateImage(ctx context.Context, prompt domain.ImagePrompt) (domain.ImagePath, error)
}

type SocialMediaAdapter interface {
	PublishImagePost(ctx context.Context, image domain.ImagePath, prompt domain.ImagePrompt) error
}
