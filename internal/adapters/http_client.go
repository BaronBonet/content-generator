package adapters

import "net/http"

// httpClient is an interface that represents an HTTP client.
// This exists, so we can mock the HTTP client, which is used in multiple adapters in our tests.
//
//go:generate mockery --name=httpClient
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
}
