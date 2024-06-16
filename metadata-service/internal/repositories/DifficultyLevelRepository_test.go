package repositories

import (
	"errors"
	"regexp"
	"testing"

	m "metadata-service/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	difficultyLevel m.DifficultyLevel = m.DifficultyLevel{
		ID:    uuid.New(),
		Level: 1,
	}
)

func TestDifficultyLevelFindAll_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDifficultyLevelRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "difficulty_levels" WHERE "difficulty_levels"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "level"}).AddRow(difficultyLevel.ID, difficultyLevel.Level))

	result, err := r.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	if result[0].Level != difficultyLevel.Level {
		t.Errorf("expected difficultyLevel name %v, but got %v", difficultyLevel.Level, result[0].Level)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestDifficultyLevelFindAll_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDifficultyLevelRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "difficulty_levels" WHERE "difficulty_levels"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("error"))

	result, err := r.FindAll()

	assert.Error(t, err)
	assert.Len(t, result, 0)
	assert.EqualError(t, err, "error")
}

func TestDifficultyLevelFindSingle_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDifficultyLevelRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "difficulty_levels" WHERE "difficulty_levels"."deleted_at" IS NULL AND "difficulty_levels"."id" = $1 ORDER BY "difficulty_levels"."id" LIMIT $2`)).
		WithArgs(
			difficultyLevel.ID,
			1,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "level"}).AddRow(difficultyLevel.ID, difficultyLevel.Level))

	result, err := r.FindSingle(difficultyLevel)

	assert.NoError(t, err)
	assert.IsType(t, m.DifficultyLevel{}, result)
	assert.Equal(t, difficultyLevel.ID, result.ID)
}

func TestDifficultyLevelFindSingle_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDifficultyLevelRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "difficulty_levels" WHERE "difficulty_levels"."deleted_at" IS NULL AND "difficulty_levels"."id" = $1 ORDER BY "difficulty_levels"."id" LIMIT $2`)).
		WithArgs(
			difficultyLevel.ID,
			1,
		).
		WillReturnError(errors.New("error"))

	_, err := r.FindSingle(difficultyLevel)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestDifficultyLevelCreate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDifficultyLevelRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "difficulty_levels" ("level","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			difficultyLevel.Level,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "level"}).AddRow(difficultyLevel.ID, difficultyLevel.Level),
		)
	mock.ExpectCommit()

	result, err := r.Create(difficultyLevel)

	assert.NoError(t, err)
	assert.IsType(t, m.DifficultyLevel{}, result)
	assert.Equal(t, difficultyLevel.ID, result.ID)
}

func TestDifficultyLevelCreate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDifficultyLevelRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "difficulty_levels" ("level","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			difficultyLevel.Level,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Create(difficultyLevel)

	assert.Error(t, err)
	assert.IsType(t, m.DifficultyLevel{}, result)
	assert.EqualError(t, err, "error")

}

func TestDifficultyLevelUpdate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDifficultyLevelRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "difficulty_levels" SET "level"=$1,"updated_at"=$2 WHERE "difficulty_levels"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs(
			difficultyLevel.Level,
			sqlmock.AnyArg(),
			difficultyLevel.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := r.Update(difficultyLevel)

	assert.NoError(t, err)
	assert.IsType(t, m.DifficultyLevel{}, result)
	assert.Equal(t, difficultyLevel.ID, result.ID)

}

func TestDifficultyLevelUpdate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDifficultyLevelRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "difficulty_levels" SET "level"=$1,"updated_at"=$2 WHERE "difficulty_levels"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs(
			difficultyLevel.Level,
			sqlmock.AnyArg(),
			difficultyLevel.ID,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Update(difficultyLevel)

	assert.Error(t, err)
	assert.IsType(t, m.DifficultyLevel{}, result)
	assert.EqualError(t, err, "error")

}

func TestDifficultyLevelDelete_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDifficultyLevelRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "difficulty_levels" SET "deleted_at"=$1 WHERE "difficulty_levels"."id" = $2 AND "difficulty_levels"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			difficultyLevel.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.Delete(difficultyLevel)

	assert.NoError(t, err)
}

func TestDifficultyLevelDelete_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDifficultyLevelRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "difficulty_levels" SET "deleted_at"=$1 WHERE "difficulty_levels"."id" = $2 AND "difficulty_levels"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			difficultyLevel.ID,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	err := r.Delete(difficultyLevel)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")

}
