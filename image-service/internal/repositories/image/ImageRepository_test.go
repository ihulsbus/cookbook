package repositories

import (
	"errors"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	m "image-service/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	image m.Image = m.Image{
		ID:         uuid.New(),
		EntityType: "recipe",
		EntityID:   uuid.New(),
		Size:       1,
		Type:       "image/jpeg",
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

func TestImageFindAll_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewImageRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE "images"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "size", "type"}).
			AddRow(
				image.ID,
				image.EntityType,
				image.EntityID,
				image.Size,
				image.Type,
			))

	result, err := r.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestImageFindAll_NotFoundErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewImageRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE "images"."deleted_at" IS NULL`)).
		WillReturnRows(&sqlmock.Rows{})

	result, err := r.FindAll()

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.Len(t, result, 0)
}

func TestImageFindAll_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewImageRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE "images"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("error"))

	result, err := r.FindAll()

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.Len(t, result, 0)
}

func TestImageFind_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewImageRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE "images"."deleted_at" IS NULL AND "images"."id" = $1 ORDER BY "images"."id" LIMIT $2`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "size", "type"}).
			AddRow(
				image.ID,
				image.EntityType,
				image.EntityID,
				image.Size,
				image.Type,
			))

	result, err := r.Find(image)

	assert.NoError(t, err)
	assert.Equal(t, image, result)
}

func TestImageFind_NotFoundErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewImageRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE "images"."deleted_at" IS NULL AND "images"."id" = $1 ORDER BY "images"."id" LIMIT $2`)).
		WillReturnRows(&sqlmock.Rows{})

	result, err := r.Find(image)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.IsType(t, m.Image{}, result)
}

func TestImageFind_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewImageRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE "images"."deleted_at" IS NULL AND "images"."id" = $1 ORDER BY "images"."id" LIMIT $2`)).
		WillReturnError(errors.New("error"))

	result, err := r.Find(image)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.Equal(t, m.Image{}, result)
}

func TestImageCreate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewImageRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "images" ("entity_type","entity_id","size","type","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
		WithArgs(
			image.EntityType,
			image.EntityID,
			image.Size,
			image.Type,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			image.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(image.ID))
	mock.ExpectCommit()
	result, err := r.Create(image)

	assert.NoError(t, err)
	assert.IsType(t, m.Image{}, result)
}

func TestImageCreate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewImageRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "images" ("entity_type","entity_id","size","type","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
		WithArgs(
			image.EntityType,
			image.EntityID,
			image.Size,
			image.Type,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			image.ID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	_, err := r.Create(image)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestImageUpdate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewImageRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "entity_type"=$1,"entity_id"=$2,"size"=$3,"type"=$4,"updated_at"=$5 WHERE "images"."deleted_at" IS NULL AND "id" = $6`)).
		WithArgs(
			image.EntityType,
			image.EntityID,
			image.Size,
			image.Type,
			sqlmock.AnyArg(),
			image.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	result, err := r.Update(image)

	assert.NoError(t, err)
	assert.IsType(t, m.Image{}, result)
}

func TestImageUpdate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewImageRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "entity_type"=$1,"entity_id"=$2,"size"=$3,"type"=$4,"updated_at"=$5 WHERE "images"."deleted_at" IS NULL AND "id" = $6`)).
		WithArgs(
			image.EntityType,
			image.EntityID,
			image.Size,
			image.Type,
			sqlmock.AnyArg(),
			image.ID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	_, err := r.Update(image)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestImageDelete_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewImageRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "deleted_at"=$1 WHERE "images"."id" = $2 AND "images"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			image.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.Delete(image)

	assert.NoError(t, err)
}

func TestImageDelete_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewImageRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "deleted_at"=$1 WHERE "images"."id" = $2 AND "images"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			image.ID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	err := r.Delete(image)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
