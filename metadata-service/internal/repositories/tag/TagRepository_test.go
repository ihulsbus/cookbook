package repositories

import (
	"errors"
	"regexp"
	"testing"

	m "metadata-service/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	co "metadata-service/internal/common/test"
)

var (
	tag m.Tag = m.Tag{
		ID:   uuid.New(),
		Name: "tag",
	}
)

func TestTagFindAll_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tags" WHERE "tags"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(tag.ID, tag.Name))

	result, err := r.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	if result[0].Name != tag.Name {
		t.Errorf("expected tag name %v, but got %v", tag.Name, result[0].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestTagFindAll_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tags" WHERE "tags"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("error"))

	result, err := r.FindAll()

	assert.Error(t, err)
	assert.Len(t, result, 0)
	assert.EqualError(t, err, "error")
}

func TestTagFindSingle_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tags" WHERE "tags"."deleted_at" IS NULL AND "tags"."id" = $1 ORDER BY "tags"."id" LIMIT $2`)).
		WithArgs(
			tag.ID,
			1,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(tag.ID, tag.Name))

	result, err := r.FindSingle(tag)

	assert.NoError(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.Equal(t, tag.ID, result.ID)
}

func TestTagFindSingle_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tags" WHERE "tags"."deleted_at" IS NULL AND "tags"."id" = $1 ORDER BY "tags"."id" LIMIT $2`)).
		WithArgs(
			tag.ID,
			1,
		).
		WillReturnError(errors.New("error"))

	_, err := r.FindSingle(tag)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestTagCreate_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tags" ("name","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			tag.Name,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).AddRow(tag.ID, tag.Name),
		)
	mock.ExpectCommit()

	result, err := r.Create(tag)

	assert.NoError(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.Equal(t, tag.ID, result.ID)
}

func TestTagCreate_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tags" ("name","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			tag.Name,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Create(tag)

	assert.Error(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.EqualError(t, err, "error")

}

func TestTagUpdate_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tags" SET "name"=$1,"updated_at"=$2 WHERE "tags"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs(
			tag.Name,
			sqlmock.AnyArg(),
			tag.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := r.Update(tag)

	assert.NoError(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.Equal(t, tag.ID, result.ID)

}

func TestTagUpdate_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tags" SET "name"=$1,"updated_at"=$2 WHERE "tags"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs(
			tag.Name,
			sqlmock.AnyArg(),
			tag.ID,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Update(tag)

	assert.Error(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.EqualError(t, err, "error")

}

func TestTagDelete_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tags" SET "deleted_at"=$1 WHERE "tags"."id" = $2 AND "tags"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			tag.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.Delete(tag)

	assert.NoError(t, err)
}

func TestTagDelete_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewTagRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tags" SET "deleted_at"=$1 WHERE "tags"."id" = $2 AND "tags"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			tag.ID,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	err := r.Delete(tag)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")

}
