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

// Query lists ingredients. Ingredients are a shared catalog (no
// owner/visibility), so there's no leak concern here (KIT-30) - unlike
// the recipe and metadata queries, this is just a straight list.
// ingredient-service has no free-text search either, so req.Query isn't
// sent. token is the caller's own bearer token ("Bearer ..."), forwarded
// as-is since the endpoint still requires authentication.
func (r *IngredientRepository) Query(ctx context.Context, req m.SearchRequest, token string) ([]m.IngredientResult, error) {
	var results []m.IngredientResult

	url := fmt.Sprintf("%s/ingredient?page=%d&limit=%d", r.serviceURL, req.Page, req.Limit)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return results, err
	}
	httpReq.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(httpReq)
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

	var page m.PaginatedResponse[m.IngredientDTO]
	if err := json.Unmarshal(body, &page); err != nil {
		return results, err
	}

	for _, ingredient := range page.Data {
		results = append(results, m.IngredientResult{
			ID:   ingredient.ID,
			Name: ingredient.Name,
		})
	}

	return results, nil
}
