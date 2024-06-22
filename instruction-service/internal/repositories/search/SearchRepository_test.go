package repositories

import (
	"errors"
	"regexp"
	"testing"

	m "instruction-service/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	co "instruction-service/internal/common/test"
)

var (
	id            uuid.UUID                  = uuid.New()
	searchRequest m.InstructionSearchRequest = m.InstructionSearchRequest{
		RecipeID: id,
	}
)

func TestSearch_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewSearchRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT recipe_instructions.instruction_id FROM "recipe_instructions" WHERE recipe_instructions.recipe_id = $1`)).
		WithArgs(
			searchRequest.RecipeID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id"}).AddRow(searchRequest.RecipeID))

	result, err := r.SearchInstruction(searchRequest)

	assert.NoError(t, err)
	assert.IsType(t, m.InstructionSearchResult{}, result)
	assert.Equal(t, searchRequest.RecipeID, result.RecipeID)
	assert.IsType(t, []uuid.UUID{}, result.InstructionIDs)
	assert.Len(t, result.InstructionIDs, 1)
	assert.Equal(t, result.InstructionIDs[0], id)
}

func TestSearch_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewSearchRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT recipe_instructions.instruction_id FROM "recipe_instructions" WHERE recipe_instructions.recipe_id = $1`)).
		WithArgs(
			searchRequest.RecipeID,
		).
		WillReturnError(errors.New("error"))

	result, err := r.SearchInstruction(searchRequest)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, m.InstructionSearchResult{}, result)
	assert.IsType(t, []uuid.UUID{}, result.InstructionIDs)
	assert.Len(t, result.InstructionIDs, 0)
}
