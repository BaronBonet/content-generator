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
		srv.logger.Error("Error when getting article", "error", err)
		return err
	}
	srv.logger.Debug("Got article", "article", article)
	imagePrompt, err := srv.llmAdapter.CreateImagePrompt(ctx, article)
	if err != nil {
		srv.logger.Error("Error when creating image prompt", "error", err)
		return err
	}
	srv.logger.Debug("Got image prompt", "imagePrompt", imagePrompt)
	image, err := srv.generationAdapter.GenerateImage(ctx, imagePrompt)
	if err != nil {
		srv.logger.Error("Error when generating image", "error", err)
		return err
	}
	srv.logger.Debug("Generated image", "image", image)
	err = srv.socialMediaAdapter.PublishImagePost(ctx, image, imagePrompt, article.Url)
	if err != nil {
		srv.logger.Error("Error when posting image", "error", err)
		return err
	}
	srv.logger.Debug("Published image")
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
