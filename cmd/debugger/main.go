package main

import (
	"context"
	"github.com/BaronBonet/content-generator/internal/adapters"
	"github.com/BaronBonet/content-generator/internal/core/service"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := adapters.NewZapLogger(zap.NewDevelopmentConfig(), true)

	NYTimesKey, exists := os.LookupEnv("NEW_YORK_TIMES_KEY")
	if !exists {
		logger.Fatal("NEW_YORK_TIMES_KEY not found")
	}
	newsAdapter := adapters.NewNYTimesNewsAdapter(NYTimesKey, http.DefaultClient)

	OpenAIKey, exists := os.LookupEnv("OPENAI_KEY")
	if !exists {
		logger.Fatal("OPENAI_KEY not found")
	}

	llmAdapter := adapters.NewChatGPTAdapter(OpenAIKey, http.DefaultClient)

	imageGenerationAdapter := adapters.NewDalleImageGenerationAdapter(OpenAIKey, http.DefaultClient)

	socialMediaAdapter, err := adapters.NewTwitterAdapterFromEnv()
	if err != nil {
		logger.Fatal("Error when creating twitter adapter", "error", err)
	}

	contentService := service.NewNewsContentService(logger, newsAdapter, llmAdapter, imageGenerationAdapter, socialMediaAdapter)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = contentService.GenerateNewsContent(ctx)
	if err != nil {
		logger.Fatal("Error when generating news content", "error", err)
	}
}
