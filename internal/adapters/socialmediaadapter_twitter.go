package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type twitterAdapter struct {
	httpClient httpClient
}

func NewTwitterAdapter(httpClient httpClient) ports.SocialMediaAdapter {
	return &twitterAdapter{
		httpClient: httpClient,
	}
}

func (t *twitterAdapter) PublishImagePost(ctx context.Context, image domain.ImagePath, prompt domain.ImagePrompt, sourceUrl string) error {
	mediaID, err := t.uploadImage(ctx, string(image))
	if err != nil {
		return err
	}

	return t.createTweet(ctx, string(prompt), mediaID)
}

func (t *twitterAdapter) createTweet(ctx context.Context, tweetText, mediaID string) error {
	tweetData := struct {
		Text        string `json:"text"`
		Attachments struct {
			MediaIDs []string `json:"media_ids"`
		} `json:"media"`
	}{
		Text: tweetText,
	}

	tweetData.Attachments.MediaIDs = append(tweetData.Attachments.MediaIDs, mediaID)

	jsonBytes, err := json.Marshal(tweetData)
	if err != nil {
		return fmt.Errorf("failed to marshal tweet data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.twitter.com/2/tweets", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		return fmt.Errorf("failed to post tweet, status code: %d", resp.StatusCode)
	}

	return nil
}

func (t *twitterAdapter) uploadImage(ctx context.Context, imageURL string) (string, error) {
	// Download the image
	resp, err := t.httpClient.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	imgData, err := ioutil.ReadAll(resp.Body)
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

	res, err := t.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
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
