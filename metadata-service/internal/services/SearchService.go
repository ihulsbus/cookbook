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
func NewSearchService(searchRepo SearchRepository) *SearchService {
	return &SearchService{
		repo: searchRepo,
	}
}

func (s SearchService) SearchMetadata(m.MetadataSearchRequestDTO) (m.MetadataSearchResultDTO, error) {
	var response m.MetadataSearchResultDTO

	return response, nil
}
