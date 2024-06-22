package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	m "metadata-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type DifficultyLevelServiceMock struct {
}

var (
	difficultyLevels []m.DifficultyLevelDTO
	difficultyLevel  m.DifficultyLevelDTO = m.DifficultyLevelDTO{
		ID:    uuid.New(),
		Level: 1,
	}
)

// ====== DifficultyLevelService ======

func (s *DifficultyLevelServiceMock) FindAll() ([]m.DifficultyLevelDTO, error) {
	switch difficultyLevel.Level {
	case 1:
		return difficultyLevels, nil
	case 2:
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (s *DifficultyLevelServiceMock) FindSingle(difficultyLevelDTO m.DifficultyLevelDTO) (m.DifficultyLevelDTO, error) {
	switch difficultyLevel.Level {
	case 1:
		return difficultyLevel, nil
	case 2:
		return m.DifficultyLevelDTO{}, errors.New("not found")
	default:
		return m.DifficultyLevelDTO{}, errors.New("error")
	}
}

func (s *DifficultyLevelServiceMock) Create(difficultyLevelDTO m.DifficultyLevelDTO) (m.DifficultyLevelDTO, error) {
	switch difficultyLevelDTO.Level {
	case 1:
		return difficultyLevel, nil
	default:
		return m.DifficultyLevelDTO{}, errors.New("error")
	}
}

func (s *DifficultyLevelServiceMock) Update(difficultyLevelDTO m.DifficultyLevelDTO) (m.DifficultyLevelDTO, error) {
	switch difficultyLevelDTO.Level {
	case 1:
		return difficultyLevel, nil
	default:
		return m.DifficultyLevelDTO{}, errors.New("error")
	}
}

func (s *DifficultyLevelServiceMock) Delete(difficultyLevelDTO m.DifficultyLevelDTO) error {
	switch difficultyLevel.Level {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

func TestDifficultyLevelGetAll_OK(t *testing.T) {
	difficultyLevels = append(difficultyLevels, difficultyLevel)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	difficultyLevel.Level = 1

	req := httptest.NewRequest("GET", "http://example.com/api/v2/difficultyLevel", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(difficultyLevels)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestDifficultyLevelGetAll_NotFound(t *testing.T) {
	difficultyLevels = append(difficultyLevels, difficultyLevel)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	difficultyLevel.Level = 2

	req := httptest.NewRequest("GET", "http://example.com/api/v2/difficultyLevel", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no difficultyLevels found"}`, string(body))
}

func TestDifficultyLevelGetAll_Error(t *testing.T) {
	difficultyLevels = append(difficultyLevels, difficultyLevel)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	difficultyLevel.Level = 3

	req := httptest.NewRequest("GET", "http://example.com/api/v2/difficultyLevel", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestDifficultyLevelGet_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/difficultyLevel/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: difficultyLevel.ID.String()},
	}

	difficultyLevel.Level = 1

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(difficultyLevel)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestDifficultyLevelGet_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/difficultyLevel/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	difficultyLevel.Level = 1

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid difficultyLevel ID"}`, string(body))
}

func TestDifficultyLevelGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/difficultyLevel/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: difficultyLevel.ID.String()},
	}

	difficultyLevel.Level = 2

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"difficultyLevel not found"}`, string(body))
}

func TestDifficultyLevelGet_FindErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/difficultyLevel/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: difficultyLevel.ID.String()},
	}

	difficultyLevel.Level = 3

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestDifficultyLevelCreate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	createDifficultyLevel := m.DifficultyLevelDTO{
		Level: 1,
	}
	reqBody, _ := json.Marshal(createDifficultyLevel)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/difficultyLevel/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assertBody, _ := json.Marshal(difficultyLevel)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, assertBody, body)
}

func TestDifficultyLevelCreate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/difficultyLevel/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unexpected JSON input"}`, string(body))
}

func TestDifficultyLevelCreate_CreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	createDifficultyLevel := m.DifficultyLevelDTO{
		Level: 3,
	}
	reqBody, _ := json.Marshal(createDifficultyLevel)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/difficultyLevel/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestDifficultyLevelUpdate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	difficultyLevel.Level = 1
	reqBody, _ := json.Marshal(difficultyLevel)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/difficultyLevel/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: difficultyLevel.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestDifficultyLevelUpdate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/difficultyLevel/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: difficultyLevel.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestDifficultyLevelUpdate_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(difficultyLevel)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/difficultyLevel/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid difficultyLevel ID"}`, string(body))
}

func TestDifficultyLevelUpdate_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	difficultyLevel.Level = 3
	reqBody, _ := json.Marshal(difficultyLevel)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/difficultyLevel/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: difficultyLevel.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestDifficultyLevelDelete_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	difficultyLevel.Level = 1

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/difficultyLevel/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: difficultyLevel.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, []byte(``), body)
}

func TestDifficultyLevelDelete_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/difficultyLevel/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"invalid difficultyLevel ID"}`), body)
}

func TestDifficultyLevelDelete_DeleteErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDifficultyLevelHandlers(&DifficultyLevelServiceMock{}, &m.LoggerInterfaceMock{})

	difficultyLevel.Level = 3

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/difficultyLevel/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: difficultyLevel.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"error"}`), body)
}
