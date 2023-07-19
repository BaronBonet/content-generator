package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/BaronBonet/go-logger/logger"
)

type service struct {
	logger              logger.Logger
	newsAdapter         ports.NewsAdapter
	llmAdapter          ports.LLMAdapter
	generationAdapter   ports.ImageGenerationAdapter
	socialMediaAdapters []ports.SocialMediaAdapter
}

func (srv *service) GenerateNewsContent(ctx context.Context) error {
	article, err := srv.newsAdapter.GetMainArticle(ctx)
	if err != nil {
		srv.logger.Error("Error when getting article", "error", err)
		return err
	}
	srv.logger.Debug("Got article", "article", article)

	prompt := fmt.Sprintf("Generate a single sentence image prompt based on the following news title and body:"+
		"\nTitle: %s"+
		"\nBody: %s"+
		"\n Do not include prompts that will be rejected by the Dalle safety system. For example mentioning dictators like Vladimir Putin."+
		"\n\n Examples of good prompts"+
		"\n- 3D render of a pink balloon dog in a violet room"+
		"\n- Illustration of a happy cat sitting on a couch in a living room with a coffee mug in its hand", article.Title, article.Body)

	imagePrompt, err := srv.llmAdapter.Chat(ctx, prompt)
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

	var wg sync.WaitGroup
	errCh := make(chan error, len(srv.socialMediaAdapters))

	for _, adapter := range srv.socialMediaAdapters {
		wg.Add(1)
		go func(adapter ports.SocialMediaAdapter) {
			srv.logger.Debug("Publishing image to social media", "adapter", adapter.GetName())
			defer wg.Done()
			if err := adapter.PublishImagePost(ctx, image, imagePrompt, srv.generationAdapter.GetGeneratorName(), article); err != nil {
				errCh <- err
			}
		}(adapter)
	}

	wg.Wait()

	close(errCh)

	for err := range errCh {
		if err != nil {
			srv.logger.Error("Error when posting image", "error", err)
		}
	}

	srv.logger.Debug("Published image to social medias")
	return nil
}

func (srv *service) CreatePrompt(ctx context.Context, prompt string) (string, error) {
	return srv.llmAdapter.Chat(ctx, prompt)
}

func (srv *service) GenerateImage(ctx context.Context, prompt string) (domain.ImagePath, error) {
	return srv.generationAdapter.GenerateImage(ctx, prompt)
}

func NewNewsContentService(
	logger logger.Logger,
	externalNewsAdapter ports.NewsAdapter,
	llmAdapter ports.LLMAdapter,
	imageGenerationAdapter ports.ImageGenerationAdapter,
	postingRepos []ports.SocialMediaAdapter,
) ports.Service {
	return &service{
		logger:              logger,
		newsAdapter:         externalNewsAdapter,
		llmAdapter:          llmAdapter,
		generationAdapter:   imageGenerationAdapter,
		socialMediaAdapters: postingRepos,
	}
}
