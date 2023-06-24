package adapters

import (
	"bytes"
	"context"
	"fmt"
	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/Davincible/goinsta"
	"image"
	"image/jpeg"
	"net/http"
	"os"
)

type instagramAdapter struct {
	username string
	password string
}

func NewInstagramSocialMediaAdapter(username string, password string) ports.SocialMediaAdapter {
	return &instagramAdapter{
		username: username,
		password: password,
	}
}

func (i *instagramAdapter) PublishImagePost(ctx context.Context, imagePath domain.ImagePath, prompt string, imageGeneratorName string, newsArticle domain.NewsArticle) error {
	insta := goinsta.New(i.username, i.password)
	err := insta.Login()
	if err != nil {
		return fmt.Errorf("failed to login to Instagram: %w", err)
	}
	resp, err := http.Get(string(imagePath))
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Convert image to JPEG
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95})
	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	reader := bytes.NewReader(buf.Bytes()) // convert buf to io.Reader

	defer resp.Body.Close()
	_, err = insta.Upload(
		&goinsta.UploadOptions{
			File:    reader,
			Caption: createInstagramCaption(prompt, imageGeneratorName, newsArticle),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}
	return nil
}

func createInstagramCaption(prompt string, imageGeneratorName string, newsArticle domain.NewsArticle) string {
	return fmt.Sprintf("AI Generated Content\n\n%s \n%s\n\n"+
		"Created by %s with the prompt:\n\n%s", newsArticle.Title, newsArticle.Url, imageGeneratorName, prompt)
}

func NewInstagramAdapterFromEnv() (ports.SocialMediaAdapter, error) {
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
		values["INSTAGRAM_USERNAME"],
		values["INSTAGRAM_PASSWORD"],
	), nil
}
