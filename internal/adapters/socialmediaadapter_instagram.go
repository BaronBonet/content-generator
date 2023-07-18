package adapters

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"

	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/BaronBonet/go-logger/logger"
	"github.com/Davincible/goinsta/v3"
)

type instagramAdapter struct {
	username string
	password string
	logger   logger.Logger
}

func NewInstagramSocialMediaAdapter(logger logger.Logger, username string, password string) ports.SocialMediaAdapter {
	return &instagramAdapter{
		username: username,
		password: password,
		logger:   logger,
	}
}

func (i *instagramAdapter) PublishImagePost(ctx context.Context, imagePath domain.ImagePath, prompt string, imageGeneratorName string, newsArticle domain.NewsArticle) error {
	i.logger.Debug("Trying to login to instagram")
	insta := goinsta.New(i.username, i.password)
	err := insta.Login()
	i.logger.Debug("Logged into instagram")
	if err != nil {
		return fmt.Errorf("failed to login to Instagram: %w", err)
	}
	resp, err := http.Get(string(imagePath))
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	i.logger.Debug("Downloaded image")

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95})
	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	reader := bytes.NewReader(buf.Bytes()) // convert buf to io.Reader
	i.logger.Debug("Converted image to io.Reader")

	defer resp.Body.Close()
	caption := createInstagramCaption(prompt, imageGeneratorName, newsArticle)
	i.logger.Debug("Uploading image with caption", "caption", caption)
	_, err = insta.Upload(
		&goinsta.UploadOptions{
			File:    reader,
			Caption: caption,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}
	return nil
}

func (i *instagramAdapter) GetName() string {
	return "Instagram"
}

func createInstagramCaption(prompt string, imageGeneratorName string, newsArticle domain.NewsArticle) string {
	return fmt.Sprintf("AI Generated Content, from the %s \n\nArticle Title: %s\n\n"+
		"Created by %s with the prompt:\n\n%s", newsArticle.Source, newsArticle.Title, imageGeneratorName, prompt)
}

func NewInstagramAdapterFromEnv(logger logger.Logger) (ports.SocialMediaAdapter, error) {
	keys := []string{"INSTAGRAM_USERNAME", "INSTAGRAM_PASSWORD"}

	values := make(map[string]string, len(keys))
	for _, key := range keys {
		value, exists := os.LookupEnv(key)
		if !exists {
			return nil, fmt.Errorf("environment variable %s not set", key)
		}
		values[key] = value
	}
	return NewInstagramSocialMediaAdapter(
		logger,
		values["INSTAGRAM_USERNAME"],
		values["INSTAGRAM_PASSWORD"],
	), nil
}
