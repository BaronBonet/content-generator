package main

import (
	"net/http"
	"os"

	"github.com/BaronBonet/content-generator/internal/adapters"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/BaronBonet/content-generator/internal/core/service"
	"github.com/BaronBonet/content-generator/internal/handlers"
	"github.com/BaronBonet/content-generator/internal/infrastructure"
	"github.com/BaronBonet/go-logger/logger"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	log := logger.NewZapLogger(true, infrastructure.Version)

	NYTimesKey, exists := os.LookupEnv("NEW_YORK_TIMES_KEY")
	if !exists {
		log.Fatal("NEW_YORK_TIMES_KEY not found")
	}
	newsAdapter := adapters.NewNYTimesNewsAdapter(NYTimesKey, http.DefaultClient)

	OpenAIKey, exists := os.LookupEnv("OPENAI_KEY")
	if !exists {
		log.Fatal("OPENAI_KEY not found")
	}

	llmAdapter := adapters.NewChatGPTAdapter(OpenAIKey, http.DefaultClient)

	imageGenerationAdapter := adapters.NewDalleImageGenerationAdapter(OpenAIKey, http.DefaultClient)

	twitterAdapter, err := adapters.NewTwitterAdapterFromEnv(log)
	if err != nil {
		log.Fatal("Error when creating twitter adapter", "error", err)
	}
	instagramAdapter, err := adapters.NewInstagramAdapterFromEnv(log)
	if err != nil {
		log.Fatal("Error when creating instagram adapter", "error", err)
	}

	contentService := service.NewNewsContentService(log, newsAdapter, llmAdapter, imageGenerationAdapter, []ports.SocialMediaAdapter{instagramAdapter, twitterAdapter})

	handler := handlers.NewAWSLambdaEventHandler(log, contentService)
	lambda.Start(handler.HandleEvent)
}
