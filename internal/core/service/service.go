package service

import (
	"context"
	"github.com/BaronBonet/content-generator/internal/core/ports"
)

type service struct {
	logger                ports.Logger
	newsAdapter           ports.NewsAdapter
	promptCreationAdapter ports.PromptCreationAdapter
	generationAdapter     ports.ImageGenerationAdapter
	socialMediaAdapter    ports.SocialMediaAdapter
}

func (srv *service) GenerateNewsContent(ctx context.Context) error {
	article, err := srv.newsAdapter.GetMainArticle(ctx)
	if err != nil {
		srv.logger.Error(ctx, "Error when getting article", "error", err)
		return err
	}
	imagePrompt, err := srv.promptCreationAdapter.CreateImagePrompt(ctx, article)
	if err != nil {
		srv.logger.Error(ctx, "Error when creating image prompt", "error", err)
		return err
	}
	localImage, err := srv.generationAdapter.GenerateImage(ctx, imagePrompt)
	if err != nil {
		srv.logger.Error(ctx, "Error when generating image", "error", err)
		return err
	}
	err = srv.socialMediaAdapter.PublishImagePost(ctx, localImage, imagePrompt)
	if err != nil {
		srv.logger.Error(ctx, "Error when posting image", "error", err)
		return err
	}

	return nil
}

func NewNewsContentService(logger ports.Logger, externalNewsRepo ports.NewsAdapter, imagePrompterRepo ports.PromptCreationAdapter, imageGenerationRepo ports.ImageGenerationAdapter, postingRepo ports.SocialMediaAdapter) ports.Service {
	return &service{
		logger:                logger,
		newsAdapter:           externalNewsRepo,
		promptCreationAdapter: imagePrompterRepo,
		generationAdapter:     imageGenerationRepo,
		socialMediaAdapter:    postingRepo,
	}
}
