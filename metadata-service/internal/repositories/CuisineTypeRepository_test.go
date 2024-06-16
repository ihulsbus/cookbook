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
	cuisineType m.CuisineType = m.CuisineType{
		ID:   uuid.New(),
		Name: "cuisineType",
	}
)

func TestCuisineTypeFindAll_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCuisineTypeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cuisine_types" WHERE "cuisine_types"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(cuisineType.ID, cuisineType.Name))

	result, err := r.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	if result[0].Name != cuisineType.Name {
		t.Errorf("expected cuisineType name %v, but got %v", cuisineType.Name, result[0].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestCuisineTypeFindAll_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCuisineTypeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cuisine_types" WHERE "cuisine_types"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("error"))

	result, err := r.FindAll()

	assert.Error(t, err)
	assert.Len(t, result, 0)
	assert.EqualError(t, err, "error")
}

func TestCuisineTypeFindSingle_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCuisineTypeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cuisine_types" WHERE "cuisine_types"."deleted_at" IS NULL AND "cuisine_types"."id" = $1 ORDER BY "cuisine_types"."id" LIMIT $2`)).
		WithArgs(
			cuisineType.ID,
			1,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(cuisineType.ID, cuisineType.Name))

	result, err := r.FindSingle(cuisineType)

	assert.NoError(t, err)
	assert.IsType(t, m.CuisineType{}, result)
	assert.Equal(t, cuisineType.ID, result.ID)
}

func TestCuisineTypeFindSingle_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCuisineTypeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cuisine_types" WHERE "cuisine_types"."deleted_at" IS NULL AND "cuisine_types"."id" = $1 ORDER BY "cuisine_types"."id" LIMIT $2`)).
		WithArgs(
			cuisineType.ID,
			1,
		).
		WillReturnError(errors.New("error"))

	_, err := r.FindSingle(cuisineType)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestCuisineTypeCreate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCuisineTypeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "cuisine_types" ("name","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			cuisineType.Name,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).AddRow(cuisineType.ID, cuisineType.Name),
		)
	mock.ExpectCommit()

	result, err := r.Create(cuisineType)

	assert.NoError(t, err)
	assert.IsType(t, m.CuisineType{}, result)
	assert.Equal(t, cuisineType.ID, result.ID)
}

func TestCuisineTypeCreate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCuisineTypeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "cuisine_types" ("name","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			cuisineType.Name,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Create(cuisineType)

	assert.Error(t, err)
	assert.IsType(t, m.CuisineType{}, result)
	assert.EqualError(t, err, "error")

}

func TestCuisineTypeUpdate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCuisineTypeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "cuisine_types" SET "name"=$1,"updated_at"=$2 WHERE "cuisine_types"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs(
			cuisineType.Name,
			sqlmock.AnyArg(),
			cuisineType.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := r.Update(cuisineType)

	assert.NoError(t, err)
	assert.IsType(t, m.CuisineType{}, result)
	assert.Equal(t, cuisineType.ID, result.ID)

}

func TestCuisineTypeUpdate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCuisineTypeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "cuisine_types" SET "name"=$1,"updated_at"=$2 WHERE "cuisine_types"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs(
			cuisineType.Name,
			sqlmock.AnyArg(),
			cuisineType.ID,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Update(cuisineType)

	assert.Error(t, err)
	assert.IsType(t, m.CuisineType{}, result)
	assert.EqualError(t, err, "error")

}

func TestCuisineTypeDelete_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCuisineTypeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "cuisine_types" SET "deleted_at"=$1 WHERE "cuisine_types"."id" = $2 AND "cuisine_types"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			cuisineType.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.Delete(cuisineType)

	assert.NoError(t, err)
}

func TestCuisineTypeDelete_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCuisineTypeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "cuisine_types" SET "deleted_at"=$1 WHERE "cuisine_types"."id" = $2 AND "cuisine_types"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			cuisineType.ID,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	err := r.Delete(cuisineType)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")

}
