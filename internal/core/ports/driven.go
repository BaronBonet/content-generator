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

//go:generate mockery --name=PromptCreationAdapter
type PromptCreationAdapter interface {
	CreateImagePrompt(ctx context.Context, article domain.NewsArticle) (domain.ImagePrompt, error)
}

//go:generate mockery --name=ImageGenerationAdapter
type ImageGenerationAdapter interface {
	GenerateImage(ctx context.Context, prompt domain.ImagePrompt) (domain.ImagePath, error)
}

//go:generate mockery --name=SocialMediaAdapter
type SocialMediaAdapter interface {
	PublishImagePost(ctx context.Context, image domain.ImagePath, prompt domain.ImagePrompt) error
}
