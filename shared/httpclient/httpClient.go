package httpClient

import (
	"net/http"
	"time"
)

// HTTPClient provides a shared HTTP client with default settings.
type HTTPClient struct {
	Client *http.Client
}

// NewHTTPClient initializes and returns an HTTPClient instance.
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Do sends an HTTP request and returns the HTTP response.
func (h *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return h.Client.Do(req)
}
