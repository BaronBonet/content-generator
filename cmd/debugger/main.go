package main

import (
	"context"
	"fmt"
	"github.com/BaronBonet/content-generator/internal/adapters"
	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"net/http"
	"os"
	"time"
)

func main() {
	// Initialize the New York Times adapter.
	key, exists := os.LookupEnv("NEW_YORK_TIMES_KEY")
	if !exists {
		return
	}

	chatKey, exists := os.LookupEnv("OPENAI_API")
	if !exists {
		return
	}

	nytAdapter := adapters.NewNYTimesNewsAdapter(key, http.DefaultClient)
	chatGPTAdapter := adapters.NewChatGPTAdapter(chatKey, http.DefaultClient)

	// Retrieve the main article using the adapter.
	mainArticle, err := getMainArticle(nytAdapter)
	if err != nil {
		return
	}
	prompt, err := getPrompt(chatGPTAdapter, mainArticle)
	println(prompt)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}

func getPrompt(adapter ports.LLMAdapter, article domain.NewsArticle) (domain.ImagePrompt, error) {
	// Set up a context with a timeout for the API request.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return adapter.CreateImagePrompt(ctx, article)
}

func getMainArticle(adapter ports.NewsAdapter) (domain.NewsArticle, error) {
	// Set up a context with a timeout for the API request.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch the main article using the adapter.
	return adapter.GetMainArticle(ctx)
}
