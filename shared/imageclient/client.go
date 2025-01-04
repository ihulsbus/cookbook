package imageclient

import (
	"fmt"
	"io"
	"net/http"
)

// ImageAPIClient interacts with the Image microservice.
type ImageAPIClient struct {
	BaseURL    string
	HTTPClient HttpClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewImageAPIClient initializes and returns a ImageAPIClient instance.
func NewImageAPIClient(baseURL string, httpClient HttpClient) *ImageAPIClient {
	return &ImageAPIClient{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
	}
}

// ImageExists checks if a image exists by its ID.
func (c *ImageAPIClient) ImageExists(imageID string) (bool, error) {
	url := fmt.Sprintf("%s/images/%s", c.BaseURL, imageID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("%+v\n", resp)

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode, body)
	}

	return true, nil
}
