package repositories

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"

	co "metadata-service/internal/common/test"
)

var (
	preparationTime m.PreparationTime = m.PreparationTime{
		ID:       uuid.New(),
		Duration: 2*time.Hour + 30*time.Minute + 15*time.Second,
	}
)

func TestPreparationTimeFindAll_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	pagination := m.NormalizePagination(1, 25)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "preparation_times"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "preparation_times" LIMIT $1`)).
		WithArgs(25).
		WillReturnRows(sqlmock.NewRows([]string{"id", "duration"}).AddRow(preparationTime.ID, preparationTime.Duration))

	data, total, err := r.FindAll(pagination)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, data, 1)

	if data[0].Duration != preparationTime.Duration {
		t.Errorf("expected preparationTime name %v, but got %v", preparationTime.Duration, data[0].Duration)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestPreparationTimeFindAll_Empty(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	pagination := m.NormalizePagination(1, 25)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "preparation_times"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "preparation_times" LIMIT $1`)).
		WithArgs(25).
		WillReturnRows(sqlmock.NewRows([]string{"id", "duration"}))

	data, total, err := r.FindAll(pagination)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, data, 0)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestPreparationTimeFindAll_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	pagination := m.NormalizePagination(1, 25)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "preparation_times"`)).
		WillReturnError(errors.New("error"))

	data, total, err := r.FindAll(pagination)

	assert.Error(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, data, 0)
	assert.EqualError(t, err, "error")
}

func TestPreparationTimeFindSingle_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "preparation_times" WHERE "preparation_times"."id" = $1 ORDER BY "preparation_times"."id" LIMIT $2`)).
		WithArgs(
			preparationTime.ID,
			1,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "duration"}).AddRow(preparationTime.ID, preparationTime.Duration))

	result, err := r.FindSingle(preparationTime)

	assert.NoError(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.Equal(t, preparationTime.ID, result.ID)
}

func TestPreparationTimeFindSingle_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "preparation_times" WHERE "preparation_times"."id" = $1 ORDER BY "preparation_times"."id" LIMIT $2`)).
		WithArgs(
			preparationTime.ID,
			1,
		).
		WillReturnError(errors.New("error"))

	_, err := r.FindSingle(preparationTime)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestPreparationTimeCreate_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "preparation_times" ("duration","id") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs(
			preparationTime.Duration,
			sqlmock.AnyArg(),
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "duration"}).AddRow(preparationTime.ID, preparationTime.Duration),
		)
	mock.ExpectCommit()

	result, err := r.Create(preparationTime)

	assert.NoError(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.Equal(t, preparationTime.ID, result.ID)
}

func TestPreparationTimeCreate_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "preparation_times" ("duration","id") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs(
			preparationTime.Duration,
			sqlmock.AnyArg(),
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Create(preparationTime)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.EqualError(t, err, "error")

}

func TestPreparationTimeUpdate_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "preparation_times" SET "duration"=$1 WHERE "id" = $2`)).
		WithArgs(
			preparationTime.Duration,
			preparationTime.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := r.Update(preparationTime)

	assert.NoError(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.Equal(t, preparationTime.ID, result.ID)

}

func TestPreparationTimeUpdate_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "preparation_times" SET "duration"=$1 WHERE "id" = $2`)).
		WithArgs(
			preparationTime.Duration,
			preparationTime.ID,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Update(preparationTime)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.EqualError(t, err, "error")

}

func TestPreparationTimeDelete_OK(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "preparation_times" WHERE "preparation_times"."id" = $1`)).
		WithArgs(
			preparationTime.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.Delete(preparationTime)

	assert.NoError(t, err)
}

func TestPreparationTimeDelete_Err(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "preparation_times" WHERE "preparation_times"."id" = $1`)).
		WithArgs(
			preparationTime.ID,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	err := r.Delete(preparationTime)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
