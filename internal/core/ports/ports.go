package ports

import (
	"context"
	"github.com/BaronBonet/content-generator/internal/core/domain"
)

type NewsAdapter interface {
	GetMainArticle(ctx context.Context) (domain.NewsArticle, error)
}

type PromptCreationAdapter interface {
	CreateImagePrompt(ctx context.Context, article domain.NewsArticle) (domain.ImagePrompt, error)
}

type ImageGenerationAdapter interface {
	GenerateImage(ctx context.Context, prompt domain.ImagePrompt) (domain.ImagePath, error)
}

type SocialMediaAdapter interface {
	PublishImagePost(ctx context.Context, localImage domain.ImagePath, imagePrompt domain.ImagePrompt) error
}

type Logger interface {
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
}
