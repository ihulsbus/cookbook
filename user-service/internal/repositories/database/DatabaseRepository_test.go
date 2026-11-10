package repositories

import (
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newMockDatabase(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	var mockDB *gorm.DB

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
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
	t, _ := time.Parse("2006-01-02 15:04", "2023-02-04 18:00")
	return t
}

func TestGetDisplayName_Found(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "display_name"}).AddRow(id, "Alice"))

	name, err := r.GetDisplayName(id)

	assert.NoError(t, err)
	assert.Equal(t, "Alice", name)
}

func TestGetDisplayName_NotFound(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "display_name"}))

	_, err := r.GetDisplayName(id)

	assert.EqualError(t, err, "user not found")
}

func TestUpsertDisplayName(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(id, "Alice", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := r.UpsertDisplayName(id, "Alice")

	assert.NoError(t, err)
}
