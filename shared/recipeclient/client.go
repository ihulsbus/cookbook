package recipeclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"
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

func (c *RecipeAPIClient) fetchPage(page, limit int) (m.PaginatedResponse[m.RecipeDTO], error) {
	url := fmt.Sprintf("%s/recipe?page=%d&limit=%d", c.BaseURL, page, limit)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return m.PaginatedResponse[m.RecipeDTO]{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return m.PaginatedResponse[m.RecipeDTO]{}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return m.PaginatedResponse[m.RecipeDTO]{}, fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode, body)
	}

	var result m.PaginatedResponse[m.RecipeDTO]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return m.PaginatedResponse[m.RecipeDTO]{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetAllRecipes fetches recipes from the recipe service.
// If no pagination is provided, all recipes are fetched across all pages.
// If a PaginationRequest is provided, only that specific page is returned.
func (c *RecipeAPIClient) GetAllRecipes(pagination ...m.PaginationRequest) (m.PaginatedResponse[m.RecipeDTO], error) {
	if len(pagination) > 0 {
		p := pagination[0]
		return c.fetchPage(p.Page, p.Limit)
	}

	const maxLimit = 100
	var allRecipes []m.RecipeDTO
	var lastMeta m.PaginationMetadata
	page := 1

	for {
		result, err := c.fetchPage(page, maxLimit)
		if err != nil {
			return m.PaginatedResponse[m.RecipeDTO]{}, err
		}
		allRecipes = append(allRecipes, result.Data...)
		lastMeta = result.Pagination
		if !result.Pagination.HasNext {
			break
		}
		page++
	}

	return m.PaginatedResponse[m.RecipeDTO]{
		Data:       allRecipes,
		Pagination: lastMeta,
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

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode, body)
	}

	return true, nil
}
