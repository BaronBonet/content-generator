package main

import (
	"github.com/BaronBonet/content-generator/internal/adapters"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/BaronBonet/content-generator/internal/core/service"
	"github.com/BaronBonet/content-generator/internal/handlers"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"net/http"
	"os"
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

	twitterAdapter, err := adapters.NewTwitterAdapterFromEnv(logger)
	if err != nil {
		logger.Fatal("Error when creating twitter adapter", "error", err)
	}
	instagramAdapter, err := adapters.NewInstagramAdapterFromEnv(logger)
	if err != nil {
		logger.Fatal("Error when creating instagram adapter", "error", err)
	}

	contentService := service.NewNewsContentService(logger, newsAdapter, llmAdapter, imageGenerationAdapter, []ports.SocialMediaAdapter{instagramAdapter, twitterAdapter})

	handler := handlers.NewAWSLambdaEventHandler(logger, contentService)
	lambda.Start(handler.HandleEvent)
}
