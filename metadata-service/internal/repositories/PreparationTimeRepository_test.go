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
	preparationTime m.PreparationTime = m.PreparationTime{
		ID:       uuid.New(),
		Duration: 1,
	}
)

func TestPreparationTimeFindAll_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "preparation_times" WHERE "preparation_times"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "duration"}).AddRow(preparationTime.ID, preparationTime.Duration))

	result, err := r.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	if result[0].Duration != preparationTime.Duration {
		t.Errorf("expected preparationTime name %v, but got %v", preparationTime.Duration, result[0].Duration)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestPreparationTimeFindAll_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "preparation_times" WHERE "preparation_times"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("error"))

	result, err := r.FindAll()

	assert.Error(t, err)
	assert.Len(t, result, 0)
	assert.EqualError(t, err, "error")
}

func TestPreparationTimeFindSingle_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "preparation_times" WHERE "preparation_times"."deleted_at" IS NULL AND "preparation_times"."id" = $1 ORDER BY "preparation_times"."id" LIMIT $2`)).
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
	db, mock := newMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "preparation_times" WHERE "preparation_times"."deleted_at" IS NULL AND "preparation_times"."id" = $1 ORDER BY "preparation_times"."id" LIMIT $2`)).
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
	db, mock := newMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "preparation_times" ("duration","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			preparationTime.Duration,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
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
	db, mock := newMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "preparation_times" ("duration","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			preparationTime.Duration,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Create(preparationTime)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.EqualError(t, err, "error")

}

func TestPreparationTimeUpdate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "preparation_times" SET "duration"=$1,"updated_at"=$2 WHERE "preparation_times"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs(
			preparationTime.Duration,
			sqlmock.AnyArg(),
			preparationTime.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := r.Update(preparationTime)

	assert.NoError(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.Equal(t, preparationTime.ID, result.ID)

}

func TestPreparationTimeUpdate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "preparation_times" SET "duration"=$1,"updated_at"=$2 WHERE "preparation_times"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs(
			preparationTime.Duration,
			sqlmock.AnyArg(),
			preparationTime.ID,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Update(preparationTime)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.EqualError(t, err, "error")

}

func TestPreparationTimeDelete_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "preparation_times" SET "deleted_at"=$1 WHERE "preparation_times"."id" = $2 AND "preparation_times"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			preparationTime.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.Delete(preparationTime)

	assert.NoError(t, err)
}

func TestPreparationTimeDelete_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewPreparationTimeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "preparation_times" SET "deleted_at"=$1 WHERE "preparation_times"."id" = $2 AND "preparation_times"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			preparationTime.ID,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	err := r.Delete(preparationTime)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")

}
