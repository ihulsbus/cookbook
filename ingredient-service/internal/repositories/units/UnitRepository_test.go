package repositories

import (
	"database/sql/driver"
	"errors"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	m "ingredient-service/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	unit m.Unit = m.Unit{
		ID:        uuid.New(),
		FullName:  "unit",
		ShortName: "u",
	}
)

func newMockDatabase(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {

	var mockDB *gorm.DB

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	sqlMockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sql mock init failed: %v", err.Error())
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlMockDB,
		PreferSimpleProtocol: true,
	})

	mockDB, err = gorm.Open(dialector, &gorm.Config{
		NowFunc: timeFunc,
		Logger:  newLogger,
	})
	if err != nil {
		t.Fatalf("gorm mock init failed: %v", err.Error())
	}

	return mockDB, mock
}

func timeFunc() time.Time {
	time, _ := time.Parse("2006-01-02 15:04", "2023-02-04 18:00")
	return time
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestUnitFindAll_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewUnitRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "units" WHERE "units"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(
				unit.ID,
				unit.FullName,
			))

	result, err := r.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestUnitFindAll_NotFoundErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewUnitRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "units" WHERE "units"."deleted_at" IS NULL`)).
		WillReturnRows(&sqlmock.Rows{})

	result, err := r.FindAll()

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.Len(t, result, 0)
}

func TestUnitFindAll_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewUnitRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "units" WHERE "units"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("error"))

	result, err := r.FindAll()

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.Len(t, result, 0)
}

func TestUnitFindSingle_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewUnitRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "units" WHERE "units"."deleted_at" IS NULL AND "units"."id" = $1 ORDER BY "units"."id" LIMIT $2`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(
				unit.ID,
				unit.FullName,
			))

	result, err := r.FindSingle(unit)

	assert.NoError(t, err)
	assert.Equal(t, unit, result)
}

func TestUnitFindSingle_NotFoundErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewUnitRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "units" WHERE "units"."deleted_at" IS NULL AND "units"."id" = $1 ORDER BY "units"."id" LIMIT $2`)).
		WillReturnRows(&sqlmock.Rows{})

	result, err := r.FindSingle(unit)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.IsType(t, m.Unit{}, result)
}

func TestUnitFindSingle_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewUnitRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "units" WHERE "units"."deleted_at" IS NULL AND "units"."id" = $1 ORDER BY "units"."id" LIMIT $2`)).
		WillReturnError(errors.New("error"))

	result, err := r.FindSingle(unit)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.Equal(t, m.Unit{}, result)
}

func TestUnitCreate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewUnitRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "units" ("full_name","short_name","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(
			unit.FullName,
			unit.ShortName,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(unit.ID))
	mock.ExpectCommit()
	result, err := r.Create(unit)

	assert.NoError(t, err)
	assert.IsType(t, m.Unit{}, result)
}

func TestUnitCreate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewUnitRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "units" ("full_name","short_name","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(
			unit.FullName,
			unit.ShortName,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	_, err := r.Create(unit)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestUnitUpdate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewUnitRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "units" SET "full_name"=$1,"short_name"=$2,"updated_at"=$3 WHERE "units"."deleted_at" IS NULL AND "id" = $4`)).
		WithArgs(
			unit.FullName,
			unit.ShortName,
			sqlmock.AnyArg(),
			unit.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	result, err := r.Update(unit)

	assert.NoError(t, err)
	assert.IsType(t, m.Unit{}, result)
}

func TestUnitUpdate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewUnitRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "units" SET "full_name"=$1,"short_name"=$2,"updated_at"=$3 WHERE "units"."deleted_at" IS NULL AND "id" = $4`)).
		WithArgs(
			unit.FullName,
			unit.ShortName,
			sqlmock.AnyArg(),
			unit.ID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	_, err := r.Update(unit)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestUnitDelete_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewUnitRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "units" SET "deleted_at"=$1 WHERE "units"."id" = $2 AND "units"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			unit.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.Delete(unit)

	assert.NoError(t, err)
}

func TestUnitDelete_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewUnitRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "units" SET "deleted_at"=$1 WHERE "units"."id" = $2 AND "units"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			unit.ID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	err := r.Delete(unit)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
