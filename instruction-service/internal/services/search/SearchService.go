package services

import (
	m "instruction-service/internal/models"
)

type SearchRepository interface {
	SearchInstruction(request m.InstructionSearchRequest) (m.InstructionSearchResult, error)
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

func (s SearchService) SearchInstruction(searchRequestDTO m.InstructionSearchRequestDTO) (m.InstructionSearchResultDTO, error) {
	var searchRequest m.InstructionSearchRequest = m.InstructionSearchRequest(searchRequestDTO)

	result, err := s.repo.SearchInstruction(searchRequest)
	if err != nil {
		return m.InstructionSearchResultDTO{}, err
	}

	return m.InstructionSearchResultDTO(result), nil
}
