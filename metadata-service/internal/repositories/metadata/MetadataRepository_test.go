package repositories

// import (
// 	"errors"
// 	"regexp"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	co "metadata-service/internal/common/test"
// 	m "metadata-service/internal/models"
// )

// var (
// 	sampleMeta = m.RecipeMetadata{
// 		RecipeID:          uuid.New(),
// 		PreparationTime:   20,
// 		ServingCount:      4,
// 		CuisineTypeID:     uuid.New(),
// 		DifficultyLevelID: uuid.New(),
// 	}
// )

// func TestFindAll_OK(t *testing.T) {
// 	db, mock := co.NewMockDatabase(t)
// 	r := NewRecipeMetadataRepository(db)

// 	rows := sqlmock.NewRows([]string{
// 		"recipe_id", "preparation_time", "serving_count", "cuisine_type_id", "difficulty_level_id",
// 	}).AddRow(
// 		sampleMeta.RecipeID,
// 		sampleMeta.PreparationTime,
// 		sampleMeta.ServingCount,
// 		sampleMeta.CuisineTypeID,
// 		sampleMeta.DifficultyLevelID,
// 	)

// 	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_metadata"`)).
// 		WillReturnRows(rows)

// 	result, err := r.FindAll()
// 	assert.NoError(t, err)
// 	assert.Len(t, result, 1)
// 	assert.Equal(t, sampleMeta.RecipeID, result[0].RecipeID)

// 	require.NoError(t, mock.ExpectationsWereMet())
// }

// func TestFindAll_Err(t *testing.T) {
// 	db, mock := co.NewMockDatabase(t)
// 	r := NewRecipeMetadataRepository(db)

// 	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_metadata"`)).
// 		WillReturnError(errors.New("error"))

// 	result, err := r.FindAll()
// 	assert.Error(t, err)
// 	assert.Nil(t, result)
// 	assert.EqualError(t, err, "error")
// }

// func TestFindSingle_OK(t *testing.T) {
// 	db, mock := co.NewMockDatabase(t)
// 	r := NewRecipeMetadataRepository(db)

// 	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_metadata" WHERE recipe_id = $1`)).
// 		WithArgs(sampleMeta.RecipeID).
// 		WillReturnRows(sqlmock.NewRows([]string{
// 			"recipe_id", "preparation_time", "serving_count", "cuisine_type_id", "difficulty_level_id",
// 		}).AddRow(
// 			sampleMeta.RecipeID,
// 			sampleMeta.PreparationTime,
// 			sampleMeta.ServingCount,
// 			sampleMeta.CuisineTypeID,
// 			sampleMeta.DifficultyLevelID,
// 		))

// 	result, err := r.FindSingle(sampleMeta)
// 	assert.NoError(t, err)
// 	assert.Equal(t, sampleMeta.RecipeID, result.RecipeID)
// }

// func TestFindSingle_Err(t *testing.T) {
// 	db, mock := co.NewMockDatabase(t)
// 	r := NewRecipeMetadataRepository(db)

// 	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_metadata" WHERE recipe_id = $1`)).
// 		WithArgs(sampleMeta.RecipeID).
// 		WillReturnError(errors.New("error"))

// 	_, err := r.FindSingle(sampleMeta)
// 	assert.Error(t, err)
// 	assert.EqualError(t, err, "error")
// }

// func TestCreate_OK(t *testing.T) {
// 	db, mock := co.NewMockDatabase(t)
// 	r := NewRecipeMetadataRepository(db)

// 	mock.ExpectBegin()
// 	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "recipe_metadata"`)).
// 		WithArgs(
// 			sampleMeta.RecipeID,
// 			sampleMeta.PreparationTime,
// 			sampleMeta.ServingCount,
// 			sampleMeta.CuisineTypeID,
// 			sampleMeta.DifficultyLevelID,
// 		).
// 		WillReturnRows(sqlmock.NewRows([]string{"recipe_id"}).AddRow(sampleMeta.RecipeID))
// 	mock.ExpectCommit()

// 	result, err := r.Create(sampleMeta)
// 	assert.NoError(t, err)
// 	assert.Equal(t, sampleMeta.RecipeID, result.RecipeID)
// }

// func TestCreate_Err(t *testing.T) {
// 	db, mock := co.NewMockDatabase(t)
// 	r := NewRecipeMetadataRepository(db)

// 	mock.ExpectBegin()
// 	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "recipe_metadata"`)).
// 		WithArgs(
// 			sampleMeta.RecipeID,
// 			sampleMeta.PreparationTime,
// 			sampleMeta.ServingCount,
// 			sampleMeta.CuisineTypeID,
// 			sampleMeta.DifficultyLevelID,
// 		).
// 		WillReturnError(errors.New("error"))
// 	mock.ExpectRollback()

// 	_, err := r.Create(sampleMeta)
// 	assert.Error(t, err)
// }

// func TestUpdate_OK(t *testing.T) {
// 	db, mock := co.NewMockDatabase(t)
// 	r := NewRecipeMetadataRepository(db)

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipe_metadata"`)).
// 		WithArgs(
// 			sampleMeta.RecipeID,
// 		).
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectCommit()

// 	result, err := r.Update(sampleMeta)
// 	assert.NoError(t, err)
// 	assert.Equal(t, sampleMeta.RecipeID, result.RecipeID)
// }

// func TestUpdate_Err(t *testing.T) {
// 	db, mock := co.NewMockDatabase(t)
// 	r := NewRecipeMetadataRepository(db)

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipe_metadata"`)).
// 		WithArgs(
// 			sampleMeta.RecipeID,
// 		).
// 		WillReturnError(errors.New("error"))
// 	mock.ExpectRollback()

// 	_, err := r.Update(sampleMeta)
// 	assert.Error(t, err)
// }

// func TestDelete_OK(t *testing.T) {
// 	db, mock := co.NewMockDatabase(t)
// 	r := NewRecipeMetadataRepository(db)

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "recipe_metadata" WHERE recipe_id = $1`)).
// 		WithArgs(sampleMeta.RecipeID).
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectCommit()

// 	err := r.Delete(sampleMeta)
// 	assert.NoError(t, err)
// }

// func TestDelete_Err(t *testing.T) {
// 	db, mock := co.NewMockDatabase(t)
// 	r := NewRecipeMetadataRepository(db)

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "recipe_metadata" WHERE recipe_id = $1`)).
// 		WithArgs(sampleMeta.RecipeID).
// 		WillReturnError(errors.New("error"))
// 	mock.ExpectRollback()

// 	err := r.Delete(sampleMeta)
// 	assert.Error(t, err)
// }
