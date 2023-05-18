package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/BaronBonet/content-generator/internal/core/domain"
	"github.com/BaronBonet/content-generator/internal/core/ports"
)

const (
	apiURL = "https://api.nytimes.com/svc/topstories/v2/home.json?api-key=%s"
)

type nyTimesAdapter struct {
	apiKey string
	client httpClient
}

// Modify your NewNYTimesNewsAdapter function to set the default client:
func NewNYTimesNewsAdapter(apiKey string, httpClient httpClient) ports.NewsAdapter {
	return &nyTimesAdapter{
		apiKey: apiKey,
		client: httpClient,
	}
}

func (n *nyTimesAdapter) GetMainArticle(ctx context.Context) (domain.NewsArticle, error) {
	url := fmt.Sprintf(apiURL, n.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.NewsArticle{}, err
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return domain.NewsArticle{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.NewsArticle{}, errors.New("failed to fetch data from New York Times API")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.NewsArticle{}, err
	}

	var apiResponse NYTApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return domain.NewsArticle{}, err
	}

	if len(apiResponse.Results) == 0 {
		return domain.NewsArticle{}, errors.New("no articles found")
	}

	mainArticle := apiResponse.Results[0]
	date, err := time.Parse(time.RFC3339, mainArticle.PublishedDate)
	if err != nil {
		return domain.NewsArticle{}, err
	}

	return domain.NewsArticle{
		Title: mainArticle.Title,
		Body:  mainArticle.Abstract,
		Date: domain.Date{
			Day:   date.Day(),
			Month: date.Month(),
			Year:  date.Year(),
		},
	}, nil
}

type NYTApiResponse struct {
	Results []NYTArticle `json:"results"`
}

type NYTArticle struct {
	Title         string `json:"title"`
	Abstract      string `json:"abstract"`
	PublishedDate string `json:"published_date"`
}
