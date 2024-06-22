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

type PreparationTimeServiceMock struct {
}

var (
	preparationTimes []m.PreparationTimeDTO
	preparationTime  m.PreparationTimeDTO = m.PreparationTimeDTO{
		ID:       uuid.New(),
		Duration: 1,
	}
)

// ====== PreparationTimeService ======

func (s *PreparationTimeServiceMock) FindAll() ([]m.PreparationTimeDTO, error) {
	switch preparationTime.Duration {
	case 1:
		return preparationTimes, nil
	case 2:
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (s *PreparationTimeServiceMock) FindSingle(preparationTimeDTO m.PreparationTimeDTO) (m.PreparationTimeDTO, error) {
	switch preparationTime.Duration {
	case 1:
		return preparationTime, nil
	case 2:
		return m.PreparationTimeDTO{}, errors.New("not found")
	default:
		return m.PreparationTimeDTO{}, errors.New("error")
	}
}

func (s *PreparationTimeServiceMock) Create(preparationTimeDTO m.PreparationTimeDTO) (m.PreparationTimeDTO, error) {
	switch preparationTimeDTO.Duration {
	case 1:
		return preparationTime, nil
	default:
		return m.PreparationTimeDTO{}, errors.New("error")
	}
}

func (s *PreparationTimeServiceMock) Update(preparationTimeDTO m.PreparationTimeDTO) (m.PreparationTimeDTO, error) {
	switch preparationTimeDTO.Duration {
	case 1:
		return preparationTime, nil
	default:
		return m.PreparationTimeDTO{}, errors.New("error")
	}
}

func (s *PreparationTimeServiceMock) Delete(preparationTimeDTO m.PreparationTimeDTO) error {
	switch preparationTime.Duration {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

func TestPreparationTimeGetAll_OK(t *testing.T) {
	preparationTimes = append(preparationTimes, preparationTime)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	preparationTime.Duration = 1

	req := httptest.NewRequest("GET", "http://example.com/api/v2/preparationTime", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(preparationTimes)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestPreparationTimeGetAll_NotFound(t *testing.T) {
	preparationTimes = append(preparationTimes, preparationTime)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	preparationTime.Duration = 2

	req := httptest.NewRequest("GET", "http://example.com/api/v2/preparationTime", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no preparationTimes found"}`, string(body))
}

func TestPreparationTimeGetAll_Error(t *testing.T) {
	preparationTimes = append(preparationTimes, preparationTime)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	preparationTime.Duration = 3

	req := httptest.NewRequest("GET", "http://example.com/api/v2/preparationTime", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestPreparationTimeGet_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/preparationTime/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: preparationTime.ID.String()},
	}

	preparationTime.Duration = 1

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(preparationTime)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestPreparationTimeGet_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/preparationTime/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	preparationTime.Duration = 1

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid preparationTime ID"}`, string(body))
}

func TestPreparationTimeGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/preparationTime/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: preparationTime.ID.String()},
	}

	preparationTime.Duration = 2

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"preparationTime not found"}`, string(body))
}

func TestPreparationTimeGet_FindErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/preparationTime/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: preparationTime.ID.String()},
	}

	preparationTime.Duration = 3

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestPreparationTimeCreate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	createPreparationTime := m.PreparationTimeDTO{
		Duration: 1,
	}
	reqBody, _ := json.Marshal(createPreparationTime)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/preparationTime/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assertBody, _ := json.Marshal(preparationTime)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, assertBody, body)
}

func TestPreparationTimeCreate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/preparationTime/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unexpected JSON input"}`, string(body))
}

func TestPreparationTimeCreate_CreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	createPreparationTime := m.PreparationTimeDTO{
		Duration: 3,
	}
	reqBody, _ := json.Marshal(createPreparationTime)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/preparationTime/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestPreparationTimeUpdate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	preparationTime.Duration = 1
	reqBody, _ := json.Marshal(preparationTime)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/preparationTime/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: preparationTime.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestPreparationTimeUpdate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/preparationTime/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: preparationTime.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestPreparationTimeUpdate_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(preparationTime)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/preparationTime/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid preparationTime ID"}`, string(body))
}

func TestPreparationTimeUpdate_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	preparationTime.Duration = 3
	reqBody, _ := json.Marshal(preparationTime)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/preparationTime/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: preparationTime.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestPreparationTimeDelete_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	preparationTime.Duration = 1

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/preparationTime/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: preparationTime.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, []byte(``), body)
}

func TestPreparationTimeDelete_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/preparationTime/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"invalid preparationTime ID"}`), body)
}

func TestPreparationTimeDelete_DeleteErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPreparationTimeHandlers(&PreparationTimeServiceMock{}, &m.LoggerInterfaceMock{})

	preparationTime.Duration = 3

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/preparationTime/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: preparationTime.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"error"}`), body)
}
