package repositories

import (
	"errors"
	"testing"

	m "metadata-service/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	tag m.Tag = m.Tag{
		TagName: "tag",
	}
)

func TestTagFindAll_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "tags" WHERE "tags"."deleted_at" IS NULL]`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := r.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestTagFindAll_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "tags" WHERE "tags"."deleted_at" IS NULL]`).
		WillReturnError(errors.New("error"))

	result, err := r.FindAll()

	assert.Error(t, err)
	assert.Len(t, result, 0)
	assert.EqualError(t, err, "error")
}

func TestTagFindSingle_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "tags" WHERE tag_id = $1 AND "tags"."deleted_at" IS NULL]`).
		WithArgs(
			1,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := r.FindSingle(uint(1))

	assert.NoError(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.Equal(t, uint(1), result.ID)
}

func TestTagFindSingle_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "tags" WHERE tag_id = $1 AND "tags"."deleted_at" IS NULL]`).
		WithArgs(
			1,
		).
		WillReturnError(errors.New("error"))

	_, err := r.FindSingle(uint(1))

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestTagCreate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`[INSERT INTO "tags" ("created_at","updated_at","deleted_at","tag_name") VALUES ($1,$2,$3,$4) RETURNING "id"]`).
		WithArgs(
			AnyTime{},
			AnyTime{},
			nil,
			"tag",
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	result, err := r.Create(tag)

	assert.NoError(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.Equal(t, uint(1), result.ID)

}

func TestTagCreate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`[INSERT INTO "tags" ("created_at","updated_at","deleted_at","tag_name") VALUES ($1,$2,$3,$4) RETURNING "id"]`).
		WithArgs(
			AnyTime{},
			AnyTime{},
			nil,
			"tag",
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Create(tag)

	assert.Error(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.EqualError(t, err, "error")

}

func TestTagUpdate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "tags" SET "updated_at"=$1,"tag_name"=$2 WHERE "id" = $3]`).
		WithArgs(
			AnyTime{},
			"tag",
			1,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updateTag := tag
	updateTag.ID = 1

	result, err := r.Update(updateTag)

	assert.NoError(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.Equal(t, uint(1), result.ID)

}

func TestTagUpdate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "tags" SET "updated_at"=$1,"tag_name"=$2 WHERE "id" = $3]`).
		WithArgs(
			AnyTime{},
			"tag",
			1,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	updateTag := tag
	updateTag.ID = 1

	result, err := r.Update(updateTag)

	assert.Error(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.EqualError(t, err, "error")

}

func TestTagDelete_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "tags" SET "deleted_at"=$1 WHERE tag_id = $2 AND "tags"."deleted_at" IS NULL]`).
		WithArgs(
			AnyTime{},
			1,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	deleteTag := tag
	deleteTag.ID = 1

	err := r.Delete(deleteTag)

	assert.NoError(t, err)
}

func TestTagDelete_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "tags" SET "updated_at"=$1,"tag_name"=$2 WHERE "id" = $3]`).
		WithArgs(
			AnyTime{},
			1,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	deleteTag := tag
	deleteTag.ID = 1

	err := r.Delete(deleteTag)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")

}
