package services

import (
	m "metadata-service/internal/models"
)

type SearchRepository interface {
	SearchMetadata(request m.MetadataSearchRequest) ([]m.MetadataSearchResult, error)
}

type SearchService struct {
	repo SearchRepository
}

// NewSearchService creates a new SearchService instance
func NewSearchService(repo SearchRepository) *SearchService {
	return &SearchService{
		repo: repo,
	}
}

func (s SearchService) SearchMetadata(searchRequestDTO m.MetadataSearchRequestDTO) ([]m.MetadataSearchResultDTO, error) {
	var searchRequest m.MetadataSearchRequest = m.MetadataSearchRequest{
		CategoryID:        &searchRequestDTO.CategoryID,
		TagID:             &searchRequestDTO.TagID,
		DifficultyLevelID: &searchRequestDTO.DifficultyLevelID,
		CuisineTypeID:     &searchRequestDTO.CuisineTypeID,
	}

	// do this separate, as otherwise both min and max will be 0 and we want nils
	if searchRequestDTO.MinPrepTime != nil {
		searchRequest.MinPrepTime = searchRequestDTO.MinPrepTime
	}
	if searchRequestDTO.MaxPrepTime != nil {
		searchRequest.MaxPrepTime = searchRequestDTO.MaxPrepTime
	}

	results, err := s.repo.SearchMetadata(searchRequest)
	if err != nil {
		return nil, err
	}

	return m.MetadataSearchResult{}.ConvertAllToDTO(results), nil
}
