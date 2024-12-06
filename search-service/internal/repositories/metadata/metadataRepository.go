package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	m "search-service/internal/models"
)

type MetadataRepository struct {
	serviceURL string
}

func NewMetadataRepository(serviceURL string) *MetadataRepository {
	return &MetadataRepository{
		serviceURL: serviceURL,
	}
}

func (r *MetadataRepository) Query(ctx context.Context, req m.SearchRequest) ([]m.MetadataSearchResult, error) {
	var results []m.MetadataSearchResult
	url := fmt.Sprintf("%s/metadata?query=%s&limit=%d&page=%d", r.serviceURL, req.Query, req.Limit, req.Page)

	resp, err := http.Get(url)
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

	err = json.Unmarshal(body, &results)
	if err != nil {
		return results, err
	}

	return results, nil
}
