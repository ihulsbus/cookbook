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

// Query calls metadata-service's unfiltered "all" listing, since that
// service has no free-text search for req.Query to map onto. The endpoint
// itself privacy-filters by the caller's token, so no filtering is needed here.
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
