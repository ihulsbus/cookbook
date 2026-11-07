package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"
)

type MetadataRepository struct {
	serviceURL string
}

func NewMetadataRepository(serviceURL string) *MetadataRepository {
	return &MetadataRepository{
		serviceURL: serviceURL,
	}
}

// Query lists recipe metadata the caller may see. metadata-service's own
// search is a structured filter (category/tag/difficulty/prep-time, see
// shared/models.MetadataSearchRequest) with no free-text mode a
// req.Query string could map onto, so this calls the unfiltered listing
// endpoint instead - every recipe's metadata the caller can see, same as
// recipe-service's own list. That endpoint is privacy-filtered against
// metadata-service's own recipe_access projection (KIT-29/30): it already
// excludes recipes the caller may not read, so there's no additional
// filtering to do here. token is the caller's own bearer token
// ("Bearer ..."), forwarded as-is - metadata-service's filter is keyed
// off whatever principal that token resolves to.
func (r *MetadataRepository) Query(ctx context.Context, req m.SearchRequest, token string) ([]m.MetadataSearchResult, error) {
	var results []m.MetadataSearchResult

	url := fmt.Sprintf("%s/metadata/search/all", r.serviceURL)

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
		return results, fmt.Errorf("metadata service responded with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}

	var dtos []m.MetadataSearchResultDTO
	if err := json.Unmarshal(body, &dtos); err != nil {
		return results, err
	}

	for _, dto := range dtos {
		results = append(results, m.MetadataSearchResult{
			RecipeID:          dto.RecipeID,
			CategoryIDs:       dto.CategoryIDs,
			TagIDs:            dto.TagIDs,
			DifficultyLevelID: dto.DifficultyLevelID,
			PreparationTimeID: dto.PreparationTimeID,
			CuisineTypeID:     dto.CuisineTypeID,
		})
	}

	return results, nil
}
