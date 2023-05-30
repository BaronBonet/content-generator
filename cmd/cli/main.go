package main

import (
	"context"
	"github.com/BaronBonet/content-generator/internal/adapters"
	"github.com/BaronBonet/content-generator/internal/core/service"
	"github.com/BaronBonet/content-generator/internal/handlers"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func main() {
	logger := adapters.NewZapLogger(zap.NewDevelopmentConfig(), true)

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

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

	socialMediaAdapter, err := adapters.NewTwitterAdapterFromEnv(logger)
	if err != nil {
		logger.Fatal("Error when creating twitter adapter", "error", err)
	}

	contentService := service.NewNewsContentService(logger, newsAdapter, llmAdapter, imageGenerationAdapter, socialMediaAdapter)
	ctx := context.Background()

	handler := handlers.NewCLIHandler(ctx, contentService)
	if err := handler.Run(os.Args); err != nil {
		logger.Fatal("Could not run CLI handler", "error", err)
	}
}