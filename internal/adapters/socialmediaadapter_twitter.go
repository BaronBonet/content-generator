package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/BaronBonet/go-logger/logger"
	"github.com/dghubble/oauth1"
)

type twitterAdapter struct {
	httpOAuthClient httpClient
	httpClient      httpClient // Used for downloading images
	logger          logger.Logger
}

type tweet struct {
	Text        string `json:"text"`
	Attachments struct {
		MediaIDs []string `json:"media_ids"`
	} `json:"media"`
}

func NewTwitterSocialMediaAdapter(httpOAuthClient httpClient, httpClient httpClient, logger logger.Logger) ports.SocialMediaAdapter {
	return &twitterAdapter{
		httpOAuthClient: httpOAuthClient,
		httpClient:      httpClient,
		logger:          logger,
	}
}

func (t *twitterAdapter) PublishImagePost(ctx context.Context, image domain.ImagePath, prompt string, imageGeneratorName string, newsArticle domain.NewsArticle) error {
	mediaID, err := t.uploadImage(ctx, string(image))
	if err != nil {
		return err
	}

	tweetID, err := t.createTweet(ctx, t.truncateString(newsArticle.Title+" "+newsArticle.Url), mediaID)
	if err != nil {
		return err
	}

	reply := fmt.Sprintf("Created by %s with the prompt:\n\n%s", imageGeneratorName, prompt)

	return t.replyToTweet(ctx, tweetID, reply)
}

func (t *twitterAdapter) GetName() string {
	return "Twitter"
}

func (t *twitterAdapter) createTweet(ctx context.Context, tweetText, mediaID string) (string, error) {

	tweetData := tweet{
		Text: tweetText,
	}
	tweetData.Attachments.MediaIDs = append(tweetData.Attachments.MediaIDs, mediaID)

	jsonBytes, err := json.Marshal(tweetData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal tweet data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.twitter.com/2/tweets", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpOAuthClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if resp.StatusCode != http.StatusCreated {
		t.logger.Error("failed to post tweet", "response body", bodyString)
		return "", fmt.Errorf("failed to post tweet, status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"data"`
	}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal tweet response: %w", err)
	}

	return response.Data.ID, nil
}

// uploadImage uploads an image to Twitter and returns the media ID, uses the v1.1 API
func (t *twitterAdapter) uploadImage(ctx context.Context, imageURL string) (string, error) {
	resp, err := t.httpClient.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("media", "media.png")
	if err != nil {
		return "", err
	}
	if _, err = fw.Write(imgData); err != nil {
		return "", err
	}
	w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://upload.twitter.com/1.1/media/upload.json", &b)

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := t.httpOAuthClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode >= 400 {
		return "", errors.New(string(body))
	}

	var data struct {
		MediaIDString string `json:"media_id_string"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}
	return data.MediaIDString, nil
}

func (t *twitterAdapter) replyToTweet(ctx context.Context, originalTweetID, replyText string) error {
	type replyTweet struct {
		Text  string `json:"text"`
		Reply struct {
			InReplyToTweetID string `json:"in_reply_to_tweet_id"`
		} `json:"reply"`
	}

	replyData := replyTweet{
		Text: t.truncateString(replyText),
	}
	replyData.Reply.InReplyToTweetID = originalTweetID

	jsonBytes, err := json.Marshal(replyData)
	if err != nil {
		return fmt.Errorf("failed to marshal reply data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.twitter.com/2/tweets", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpOAuthClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.logger.Error("failed to post tweet reply", "response body", string(bodyBytes))
		return fmt.Errorf("failed to post reply tweet, status code: %d", resp.StatusCode)
	}

	return nil
}

// truncateString breaks a string up into a chunk of size up to 280.
func (t *twitterAdapter) truncateString(s string) string {
	runeStr := []rune(s) // Convert to runes for proper handling of special characters

	if len(runeStr) > 280 {
		t.logger.Warn("Tweet was truncated to 280 characters", "full tweet", s)
		runeStr = runeStr[:280]
	}

	return string(runeStr)
}

// NewTwitterAdapterFromEnv is a helper function to create a TwitterAdapter from environment variables
func NewTwitterAdapterFromEnv(logger logger.Logger) (ports.SocialMediaAdapter, error) {
	keys := []string{"TWITTER_API_KEY", "TWITTER_API_KEY_SECRET", "TWITTER_ACCESS_TOKEN", "TWITTER_ACCESS_TOKEN_SECRET"}

	values := make(map[string]string, len(keys))
	for _, key := range keys {
		value, exists := os.LookupEnv(key)
		if !exists {
			return nil, fmt.Errorf("environment variable %s not set", key)
		}
		values[key] = value
	}

	config := oauth1.NewConfig(values["TWITTER_API_KEY"], values["TWITTER_API_KEY_SECRET"])
	token := oauth1.NewToken(values["TWITTER_ACCESS_TOKEN"], values["TWITTER_ACCESS_TOKEN_SECRET"])

	return NewTwitterSocialMediaAdapter(config.Client(oauth1.NoContext, token), http.DefaultClient, logger), nil
}
