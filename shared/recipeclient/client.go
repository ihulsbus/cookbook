package recipeclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	hc "../httpclient"
)

// RecipeAPIClient interacts with the Recipe microservice.
type RecipeAPIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewRecipeAPIClient initializes and returns a RecipeAPIClient instance.
func NewRecipeAPIClient(baseURL string, timeout time.Duration) *RecipeAPIClient {
	return &RecipeAPIClient{
		BaseURL:    baseURL,
		HTTPClient: hc.NewHTTPClient(timeout),
	}
}

// RecipeExists checks if a recipe exists by its ID.
func (c *RecipeAPIClient) RecipeExists(recipeID string) (bool, error) {
	url := fmt.Sprintf("%s/recipes/%s", c.BaseURL, recipeID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode, body)
	}

	return true, nil
}
