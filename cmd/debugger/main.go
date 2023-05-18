package main

import (
	"context"
	"github.com/BaronBonet/content-generator/internal/adapters"
	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/dghubble/oauth1"
	"os"
	"time"
)

func main() {

	twitterApiKey, exists := os.LookupEnv("TWITTER_API_KEY")
	if !exists {
		return
	}
	twitterApiSecret := os.Getenv("TWITTER_API_KEY_SECRET")
	twitterApiToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	twitterAccessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	config := oauth1.NewConfig(twitterApiKey, twitterApiSecret)
	token := oauth1.NewToken(twitterApiToken, twitterAccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	ta := adapters.NewTwitterAdapter(httpClient)
	err := ta.PublishImagePost(context.Background(), "https://cdn.ericcbonet.com/chatgpt-writes-test-part-2.png", "again", "sourceUrl")
	if err != nil {
		return
	}

	// Initialize the New York Times adapter.
	//key, exists := os.LookupEnv("NEW_YORK_TIMES_KEY")
	if !exists {
		return
	}
	//
	//chatKey, exists := os.LookupEnv("OPENAI_API")
	//if !exists {
	//	return
	//}
	//
	//nytAdapter := adapters.NewNYTimesNewsAdapter(key, http.DefaultClient)
	//chatGPTAdapter := adapters.NewChatGPTAdapter(chatKey, http.DefaultClient)
	//
	//// Retrieve the main article using the adapter.
	//mainArticle, err := getMainArticle(nytAdapter)
	//if err != nil {
	//	return
	//}
	//prompt, err := getPrompt(chatGPTAdapter, mainArticle)
	//// Prompt: A cartoon depiction of the Group of 7 leaders playing a game of ‘Russian Roulette’ with oil barrels, with Russia anxiously looking on in the background.
	//println(prompt)
	//if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	return
	//}
	//dalleAdapter := adapters.NewDalleAdapter(chatKey, http.DefaultClient)
	//imagePath, err := dalleAdapter.GenerateImage(context.Background(), prompt)
	//println(imagePath)
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
