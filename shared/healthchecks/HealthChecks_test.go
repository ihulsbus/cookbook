package health

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MockLogger struct{}

func (m *MockLogger) Debugf(format string, args ...interface{}) {}
func (m *MockLogger) Warnf(format string, args ...interface{})  {}
func (m *MockLogger) Errorf(format string, args ...interface{}) {}

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return db, mock, sqlDB
}

func TestLiveness(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	handler := NewHealthHandlers(db, &MockLogger{})

	router := gin.New()
	router.GET("/health/live", handler.Liveness)

	req, _ := http.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "UP", response.Status)
	assert.NotEmpty(t, response.Timestamp)
}

func TestReadiness_DatabaseUp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	// Mock successful database ping
	mock.ExpectPing()

	handler := NewHealthHandlers(db, &MockLogger{})

	router := gin.New()
	router.GET("/health/ready", handler.Readiness)

	req, _ := http.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "UP", response.Status)
	assert.Equal(t, "UP", response.Checks["database"])
}

//func TestReadiness_DatabaseDown(t *testing.T) {
//	gin.SetMode(gin.TestMode)
//	db, mock, sqlDB := setupMockDB(t)
//	defer sqlDB.Close()
//
//
//	// Mock failed database ping
//	mock.ExpectPing().WillReturnError(sql.ErrConnDone)
//
//	handler := NewHealthHandlers(db, &MockLogger{})
//
//	router := gin.New()
//	router.GET("/health/ready", handler.Readiness)
//
//	req, _ := http.NewRequest("GET", "/health/ready", nil)
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
//
//	var response Response
//	err := json.Unmarshal(w.Body.Bytes(), &response)
//	assert.NoError(t, err)
//	assert.Equal(t, "DOWN", response.Status)
//	assert.Contains(t, response.Checks["database"], "DOWN")
//}

func TestStartup_DatabaseUp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	// Mock successful database ping
	mock.ExpectPing()

	handler := NewHealthHandlers(db, &MockLogger{})

	router := gin.New()
	router.GET("/health/startup", handler.Startup)

	req, _ := http.NewRequest("GET", "/health/startup", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "UP", response.Status)
	assert.Equal(t, "UP", response.Checks["database"])
}
