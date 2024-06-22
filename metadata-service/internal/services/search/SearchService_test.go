package services

import (
	"errors"
	m "metadata-service/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	minTime int       = 1
	maxTime int       = 2
	id      uuid.UUID = uuid.New()

	searchRequest m.MetadataSearchRequestDTO = m.MetadataSearchRequestDTO{
		RecipeID:          id,
		CategoryID:        id,
		TagID:             id,
		DifficultyLevelID: id,
		CuisineTypeID:     id,
		MinPrepTime:       &minTime,
		MaxPrepTime:       &maxTime,
	}
)

type searchRepositoryMock struct{}

func (*searchRepositoryMock) SearchMetadata(request m.MetadataSearchRequest) ([]m.MetadataSearchResult, error) {
	switch *request.MinPrepTime {
	case 1:
		var response []m.MetadataSearchResult
		response = append(response, m.MetadataSearchResult{
			RecipeID:          id,
			CategoryIDs:       []uuid.UUID{id},
			TagIDs:            []uuid.UUID{id},
			DifficultyLevelID: id,
			CuisineTypeID:     id,
			PreparationTimeID: id,
		})
		return response, nil
	default:
		return nil, errors.New("error")
	}
}

// ======================================================================

func TestSearchMetadata_OK(t *testing.T) {
	s := NewSearchService(&searchRepositoryMock{})

	result, err := s.SearchMetadata(searchRequest)

	assert.NoError(t, err)
	assert.IsType(t, []m.MetadataSearchResultDTO{}, result)
}

func TestSearchMetadata_Error(t *testing.T) {
	s := NewSearchService(&searchRepositoryMock{})

	minTime = 2

	_, err := s.SearchMetadata(searchRequest)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
