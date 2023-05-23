package ports

import (
	"context"
	"github.com/BaronBonet/content-generator/internal/core/domain"
)

type Service interface {
	GenerateNewsContent(ctx context.Context) error
	CreatePrompt(ctx context.Context, prompt string) (string, error)
	GenerateImage(ctx context.Context, prompt string) (domain.ImagePath, error)
}
