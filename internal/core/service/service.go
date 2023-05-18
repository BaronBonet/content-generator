package service

import (
	"context"
	"github.com/BaronBonet/content-generator/internal/core/ports"
)

type service struct {
	logger             ports.Logger
	newsAdapter        ports.NewsAdapter
	llmAdapter         ports.LLMAdapter
	generationAdapter  ports.ImageGenerationAdapter
	socialMediaAdapter ports.SocialMediaAdapter
}

func (srv *service) GenerateNewsContent(ctx context.Context) error {
	article, err := srv.newsAdapter.GetMainArticle(ctx)
	if err != nil {
		srv.logger.Error(ctx, "Error when getting article", "error", err)
		return err
	}
	imagePrompt, err := srv.llmAdapter.CreateImagePrompt(ctx, article)
	if err != nil {
		srv.logger.Error(ctx, "Error when creating image prompt", "error", err)
		return err
	}
	localImage, err := srv.generationAdapter.GenerateImage(ctx, imagePrompt)
	if err != nil {
		srv.logger.Error(ctx, "Error when generating image", "error", err)
		return err
	}
	err = srv.socialMediaAdapter.PublishImagePost(ctx, localImage, imagePrompt, article.Url)
	if err != nil {
		srv.logger.Error(ctx, "Error when posting image", "error", err)
		return err
	}

	return nil
}

func NewNewsContentService(
	logger ports.Logger,
	externalNewsAdapter ports.NewsAdapter,
	llmAdapter ports.LLMAdapter,
	imageGenerationAdapter ports.ImageGenerationAdapter,
	postingRepo ports.SocialMediaAdapter,
) ports.Service {
	return &service{
		logger:             logger,
		newsAdapter:        externalNewsAdapter,
		llmAdapter:         llmAdapter,
		generationAdapter:  imageGenerationAdapter,
		socialMediaAdapter: postingRepo,
	}
}
