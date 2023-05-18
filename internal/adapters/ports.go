package adapters

import "net/http"

// client httpClient
//
//go:generate mockery --name=httpClient
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
}
