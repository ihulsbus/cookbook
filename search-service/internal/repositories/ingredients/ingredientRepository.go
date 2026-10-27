package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"
)

type IngredientRepository struct {
	serviceURL string
}

func NewIngredientRepository(serviceURL string) *IngredientRepository {
	return &IngredientRepository{
		serviceURL: serviceURL,
	}
}

func (r *IngredientRepository) Query(ctx context.Context, req m.SearchRequest) ([]m.IngredientResult, error) {
	var results []m.IngredientResult
	url := fmt.Sprintf("%s/ingredients?query=%s&limit=%d&page=%d", r.serviceURL, req.Query, req.Limit, req.Page)

	resp, err := http.Get(url)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return results, fmt.Errorf("ingredient service responded with status: %d", resp.StatusCode)
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
