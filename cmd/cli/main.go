package main

import (
	"context"
	"net/http"
	"os"

	"github.com/BaronBonet/content-generator/internal/adapters"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/BaronBonet/content-generator/internal/core/service"
	"github.com/BaronBonet/content-generator/internal/handlers"
	"github.com/BaronBonet/go-logger/logger"
	"github.com/joho/godotenv"
)

func main() {
	logger := logger.NewSlogLogger()

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

	twitterAdapter, err := adapters.NewTwitterAdapterFromEnv(logger)
	if err != nil {
		logger.Fatal("Error when creating twitter adapter", "error", err)
	}
	instagramAdapter, err := adapters.NewInstagramAdapterFromEnv(logger)

	contentService := service.NewNewsContentService(logger, newsAdapter, llmAdapter, imageGenerationAdapter, []ports.SocialMediaAdapter{instagramAdapter, twitterAdapter})
	ctx := context.Background()

	handler := handlers.NewCLIHandler(ctx, contentService, logger)
	if err := handler.Run(os.Args); err != nil {
		logger.Fatal("Could not run CLI handler", "error", err)
	}
}
