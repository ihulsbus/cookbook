package repositories

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"

	co "metadata-service/internal/common/test"
)

var (
	category m.Category = m.Category{
		ID:   uuid.New(),
		Name: "category",
	}
)

// ========================================================================================================

func TestCategoryFindAll_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewCategoryRepository(db)

	pagination := m.NormalizePagination(1, 25)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "categories" WHERE "categories"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."deleted_at" IS NULL LIMIT $1 OFFSET $2`)).
		WithArgs(25, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(category.ID, category.Name))

	data, total, err := r.FindAll(pagination)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, data, 1)

	if data[0].Name != category.Name {
		t.Errorf("expected category name %v, but got %v", category.Name, data[0].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestCategoryFindAll_Empty(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewCategoryRepository(db)

	pagination := m.NormalizePagination(1, 25)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "categories" WHERE "categories"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."deleted_at" IS NULL LIMIT $1 OFFSET $2`)).
		WithArgs(25, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	data, total, err := r.FindAll(pagination)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, data, 0)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestCategoryFindAll_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewCategoryRepository(db)

	pagination := m.NormalizePagination(1, 25)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "categories" WHERE "categories"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("error"))

	data, total, err := r.FindAll(pagination)

	assert.Error(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, data, 0)
	assert.EqualError(t, err, "error")
}

func TestCategoryFindSingle_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."deleted_at" IS NULL AND "categories"."id" = $1 ORDER BY "categories"."id" LIMIT $2`)).
		WithArgs(
			category.ID,
			1,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(category.ID, category.Name))

	result, err := r.FindSingle(category)

	assert.NoError(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.Equal(t, category.ID, result.ID)
}

func TestCategoryFindSingle_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."deleted_at" IS NULL AND "categories"."id" = $1 ORDER BY "categories"."id" LIMIT $2`)).
		WithArgs(
			category.ID,
			1,
		).
		WillReturnError(errors.New("error"))

	_, err := r.FindSingle(category)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestCategoryCreate_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories" ("name","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			category.Name,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(category.ID, category.Name))
	mock.ExpectCommit()

	result, err := r.Create(category)

	assert.NoError(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.Equal(t, category.ID, result.ID)

}

func TestCategoryCreate_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories" ("name","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			category.Name,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Create(category)

	assert.Error(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.EqualError(t, err, "error")

}

func TestCategoryUpdate_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories" SET "name"=$1,"updated_at"=$2 WHERE "categories"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs(
			category.Name,
			sqlmock.AnyArg(),
			category.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := r.Update(category)

	assert.NoError(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.Equal(t, category.ID, result.ID)

}

func TestCategoryUpdate_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories" SET "name"=$1,"updated_at"=$2 WHERE "categories"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs(
			category.Name,
			sqlmock.AnyArg(),
			category.ID,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Update(category)

	assert.Error(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.EqualError(t, err, "error")

}

func TestCategoryDelete_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories" SET "deleted_at"=$1 WHERE "categories"."id" = $2 AND "categories"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			category.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.Delete(category)

	assert.NoError(t, err)
}

func TestCategoryDelete_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories" SET "deleted_at"=$1 WHERE "categories"."id" = $2 AND "categories"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			category.ID,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	err := r.Delete(category)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")

}
