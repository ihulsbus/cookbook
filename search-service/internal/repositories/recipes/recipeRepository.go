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

// Query lists recipes visible to the caller (KIT-30). recipe-service has
// no free-text search of its own, so req.Query isn't sent - there's
// nothing on the other end to filter by it; what's returned is exactly
// what GET /recipe would return for this caller: their own recipes plus
// every public one. token is the caller's own bearer token
// ("Bearer ..."), forwarded as-is so recipe-service's own owner/
// visibility filter (KIT-25) - not this service - decides what's
// visible; nothing here must do its own filtering on top.
func (r *RecipeRepository) Query(ctx context.Context, req m.SearchRequest, token string) ([]m.RecipeResult, error) {
	var results []m.RecipeResult

	url := fmt.Sprintf("%s/recipe?page=%d&limit=%d", r.serviceURL, req.Page, req.Limit)

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
		return results, fmt.Errorf("recipe service responded with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}

	var page m.PaginatedResponse[m.RecipeDTO]
	if err := json.Unmarshal(body, &page); err != nil {
		return results, err
	}

	for _, recipe := range page.Data {
		results = append(results, m.RecipeResult{
			ID:           recipe.ID,
			Name:         recipe.Name,
			Description:  recipe.Description,
			ServingCount: recipe.ServingCount,
		})
	}

	return results, nil
}
