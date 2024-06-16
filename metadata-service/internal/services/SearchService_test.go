package services

import (
	m "metadata-service/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	searchRequest m.MetadataSearchRequestDTO = m.MetadataSearchRequestDTO{
		RecipeID:        uuid.New(),
		CategoryID:      uuid.New(),
		TagID:           uuid.New(),
		DifficultyLevel: 1,
		CuisineType:     "",
		MinPrepTime:     1,
		MaxPrepTime:     2,
	}
)

type searchRepositoryMock struct{}

func (*searchRepositoryMock) SearchMetadata(request m.MetadataSearchRequest) ([]m.MetadataSearchResult, error) {
	var response []m.MetadataSearchResult

	return response, nil
}

// ======================================================================

func TestSearchMetadata_OK(t *testing.T) {
	s := NewSearchService(&searchRepositoryMock{})

	result, err := s.SearchMetadata(searchRequest)

	assert.NoError(t, err)
	assert.IsType(t, m.MetadataSearchResultDTO{}, result)
}
