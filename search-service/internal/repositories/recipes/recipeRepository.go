package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"
)

type RecipeRepository struct {
	serviceURL string
}

func NewRecipeRepository(serviceURL string) *RecipeRepository {
	return &RecipeRepository{
		serviceURL: serviceURL,
	}
}

func (r *RecipeRepository) Query(ctx context.Context, req m.SearchRequest) ([]m.RecipeResult, error) {
	var results []m.RecipeResult
	url := fmt.Sprintf("%s/recipes?query=%s&limit=%d&page=%d", r.serviceURL, req.Query, req.Limit, req.Page)

	resp, err := http.Get(url)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return results, fmt.Errorf("recipe service responded with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}

	err = json.Unmarshal(body, &results)
	if err != nil {
		return results, err
	}

	return results, nil
}
