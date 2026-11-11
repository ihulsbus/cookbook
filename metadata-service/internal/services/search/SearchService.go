package services

import (
	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
)

type SearchRepository interface {
	GetAllRecipeMetadata(requesterID uuid.UUID) (*[]m.MetadataSearchResult, error)
	SearchMetadata(request m.MetadataSearchRequest, requesterID uuid.UUID) ([]m.MetadataSearchResult, error)
}

type SearchService struct {
	repo SearchRepository
}

func NewSearchService(repo SearchRepository) *SearchService {
	return &SearchService{
		repo: repo,
	}
}

func (s SearchService) GetAllRecipeMetadata(requesterID uuid.UUID) (*[]m.MetadataSearchResultDTO, error) {
	results, err := s.repo.GetAllRecipeMetadata(requesterID)
	if err != nil {
		return nil, err
	}

	resultsDTO := m.MetadataSearchResult{}.ConvertAllToDTO(*results)
	return &resultsDTO, nil
}

func (s SearchService) SearchMetadata(searchRequestDTO m.MetadataSearchRequestDTO, requesterID uuid.UUID) ([]m.MetadataSearchResultDTO, error) {
	var searchRequest m.MetadataSearchRequest = m.MetadataSearchRequest{
		CategoryID:        &searchRequestDTO.CategoryID,
		TagID:             &searchRequestDTO.TagID,
		DifficultyLevelID: &searchRequestDTO.DifficultyLevelID,
		CuisineTypeID:     &searchRequestDTO.CuisineTypeID,
	}

	// Set separately: a direct assignment would turn an absent value into 0 instead of nil.
	if searchRequestDTO.MinPrepTime != nil {
		searchRequest.MinPrepTime = searchRequestDTO.MinPrepTime
	}
	if searchRequestDTO.MaxPrepTime != nil {
		searchRequest.MaxPrepTime = searchRequestDTO.MaxPrepTime
	}

	results, err := s.repo.SearchMetadata(searchRequest, requesterID)
	if err != nil {
		return nil, err
	}

	return m.MetadataSearchResult{}.ConvertAllToDTO(results), nil
}
