package repositories

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	m "github.com/ihulsbus/cookbook/shared/models"
	co "metadata-service/internal/common/test"
)

var (
	sampleMeta = m.RecipeMetadata{
		RecipeID:          uuid.New(),
		PreparationTime:   20,
		ServingCount:      4,
		CuisineTypeID:     uuid.New(),
		DifficultyLevelID: uuid.New(),
	}
)

func TestMetadataFindAll_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewRecipeMetadataRepository(db)

	pagination := m.NormalizePagination(1, 25)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "recipe_metadata"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_metadata" LIMIT $1 OFFSET $2`)).
		WithArgs(25, 0).
		WillReturnRows(sqlmock.NewRows([]string{
			"recipe_id", "preparation_time", "serving_count", "cuisine_type_id", "difficulty_level_id",
		}).AddRow(
			sampleMeta.RecipeID,
			sampleMeta.PreparationTime,
			sampleMeta.ServingCount,
			sampleMeta.CuisineTypeID,
			sampleMeta.DifficultyLevelID,
		))

	data, total, err := r.FindAll(pagination)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, data, 1)
	assert.Equal(t, sampleMeta.RecipeID, data[0].RecipeID)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMetadataFindAll_Empty(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewRecipeMetadataRepository(db)

	pagination := m.NormalizePagination(1, 25)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "recipe_metadata"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_metadata" LIMIT $1 OFFSET $2`)).
		WithArgs(25, 0).
		WillReturnRows(sqlmock.NewRows([]string{
			"recipe_id", "preparation_time", "serving_count", "cuisine_type_id", "difficulty_level_id",
		}))

	data, total, err := r.FindAll(pagination)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, data, 0)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMetadataFindAll_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewRecipeMetadataRepository(db)

	pagination := m.NormalizePagination(1, 25)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "recipe_metadata"`)).
		WillReturnError(errors.New("error"))

	data, total, err := r.FindAll(pagination)
	assert.Error(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, data, 0)
	assert.EqualError(t, err, "error")
}
