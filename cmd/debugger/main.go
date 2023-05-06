package main

import (
	"context"
	"fmt"
	"github.com/BaronBonet/content-generator/internal/adapters"
	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"os"
	"time"
)

func main() {
	// Initialize the New York Times adapter.
	key, exists := os.LookupEnv("NEW_YORK_TIMES_KEY")
	if !exists {
		return
	}

	nytAdapter := adapters.NewNYTimesNewsAdapter(key)

	// Retrieve the main article using the adapter.
	mainArticle, err := getMainArticle(nytAdapter)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(mainArticle)
}

func getMainArticle(adapter ports.NewsAdapter) (domain.NewsArticle, error) {
	// Set up a context with a timeout for the API request.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch the main article using the adapter.
	return adapter.GetMainArticle(ctx)
}
