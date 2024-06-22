package services

import (
	"errors"
	m "instruction-service/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	id uuid.UUID = uuid.New()

	searchRequest m.InstructionSearchRequestDTO = m.InstructionSearchRequestDTO{
		RecipeID: id,
	}

	switchCheck string
)

type searchRepositoryMock struct{}

func (*searchRepositoryMock) SearchInstruction(request m.InstructionSearchRequest) (m.InstructionSearchResult, error) {
	switch switchCheck {
	case "search":
		return m.InstructionSearchResult{
			RecipeID:       id,
			InstructionIDs: []uuid.UUID{id},
		}, nil
	default:
		return m.InstructionSearchResult{}, errors.New("error")
	}
}

// ======================================================================

func TestSearchInstruction_OK(t *testing.T) {
	s := NewSearchService(&searchRepositoryMock{})

	switchCheck = "search"

	result, err := s.SearchInstruction(searchRequest)

	assert.NoError(t, err)
	assert.IsType(t, m.InstructionSearchResultDTO{}, result)
}

func TestSearchInstruction_Error(t *testing.T) {
	s := NewSearchService(&searchRepositoryMock{})

	switchCheck = "error"

	_, err := s.SearchInstruction(searchRequest)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
