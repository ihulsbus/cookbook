package recipeclient

import (
	"fmt"
	"io"
	"net/http"
)

// RecipeAPIClient interacts with the Recipe microservice.
type RecipeAPIClient struct {
	BaseURL    string
	HTTPClient HttpClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewRecipeAPIClient initializes and returns a RecipeAPIClient instance.
func NewRecipeAPIClient(baseURL string, httpClient HttpClient) (*RecipeAPIClient, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL cannot be empty")
	}
	if httpClient == nil {
		return nil, fmt.Errorf("httpClient cannot be nil")
	}
	return &RecipeAPIClient{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
	}, nil
}

// RecipeExists checks if a recipe exists by its ID.
func (c *RecipeAPIClient) RecipeExists(recipeID string) (bool, error) {
	url := fmt.Sprintf("%s/recipe/%s", c.BaseURL, recipeID)
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
