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

type CuisineTypeServiceMock struct {
}

var (
	cuisineTypes []m.CuisineTypeDTO
	cuisineType  m.CuisineTypeDTO = m.CuisineTypeDTO{
		ID:   uuid.New(),
		Name: "cuisineType",
	}
)

// ====== CuisineTypeService ======

func (s *CuisineTypeServiceMock) FindAll() ([]m.CuisineTypeDTO, error) {
	switch cuisineType.Name {
	case "findall":
		return cuisineTypes, nil
	case "notfound":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (s *CuisineTypeServiceMock) FindSingle(cuisineTypeDTO m.CuisineTypeDTO) (m.CuisineTypeDTO, error) {
	switch cuisineType.Name {
	case "find":
		return cuisineType, nil
	case "notfound":
		return m.CuisineTypeDTO{}, errors.New("not found")
	default:
		return m.CuisineTypeDTO{}, errors.New("error")
	}
}

func (s *CuisineTypeServiceMock) Create(cuisineTypeDTO m.CuisineTypeDTO) (m.CuisineTypeDTO, error) {
	switch cuisineTypeDTO.Name {
	case "create":
		return cuisineType, nil
	default:
		return m.CuisineTypeDTO{}, errors.New("error")
	}
}

func (s *CuisineTypeServiceMock) Update(cuisineTypeDTO m.CuisineTypeDTO) (m.CuisineTypeDTO, error) {
	switch cuisineTypeDTO.Name {
	case "update":
		return cuisineType, nil
	default:
		return m.CuisineTypeDTO{}, errors.New("error")
	}
}

func (s *CuisineTypeServiceMock) Delete(cuisineTypeDTO m.CuisineTypeDTO) error {
	switch cuisineType.Name {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

func TestCuisineTypeGetAll_OK(t *testing.T) {
	cuisineTypes = append(cuisineTypes, cuisineType)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	cuisineType.Name = "findall"

	req := httptest.NewRequest("GET", "http://example.com/api/v1/cuisineType", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(cuisineTypes)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestCuisineTypeGetAll_NotFound(t *testing.T) {
	cuisineTypes = append(cuisineTypes, cuisineType)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	cuisineType.Name = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v1/cuisineType", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no cuisineTypes found"}`, string(body))
}

func TestCuisineTypeGetAll_Error(t *testing.T) {
	cuisineTypes = append(cuisineTypes, cuisineType)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	cuisineType.Name = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v1/cuisineType", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestCuisineTypeGet_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/cuisineType/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: cuisineType.ID.String()},
	}

	cuisineType.Name = "find"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(cuisineType)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestCuisineTypeGet_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/cuisineType/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	cuisineType.Name = "find"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid cuisineType ID"}`, string(body))
}

func TestCuisineTypeGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/cuisineType/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: cuisineType.ID.String()},
	}

	cuisineType.Name = "notfound"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"cuisineType not found"}`, string(body))
}

func TestCuisineTypeGet_FindErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/cuisineType/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: cuisineType.ID.String()},
	}

	cuisineType.Name = "finderr"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestCuisineTypeCreate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	createCuisineType := m.CuisineTypeDTO{
		Name: "create",
	}
	reqBody, _ := json.Marshal(createCuisineType)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/cuisineType/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assertBody, _ := json.Marshal(cuisineType)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, assertBody, body)
}

func TestCuisineTypeCreate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v1/cuisineType/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unexpected JSON input"}`, string(body))
}

func TestCuisineTypeCreate_CreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	createCuisineType := m.CuisineTypeDTO{
		Name: "createerr",
	}
	reqBody, _ := json.Marshal(createCuisineType)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/cuisineType/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestCuisineTypeUpdate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	cuisineType.Name = "update"
	reqBody, _ := json.Marshal(cuisineType)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/cuisineType/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: cuisineType.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestCuisineTypeUpdate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/cuisineType/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: cuisineType.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestCuisineTypeUpdate_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(cuisineType)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/cuisineType/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid cuisineType ID"}`, string(body))
}

func TestCuisineTypeUpdate_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	cuisineType.Name = "updatefail"
	reqBody, _ := json.Marshal(cuisineType)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/cuisineType/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: cuisineType.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestCuisineTypeDelete_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	cuisineType.Name = "delete"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/cuisineType/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: cuisineType.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, []byte(``), body)
}

func TestCuisineTypeDelete_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/cuisineType/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"invalid cuisineType ID"}`), body)
}

func TestCuisineTypeDelete_DeleteErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCuisineTypeHandlers(&CuisineTypeServiceMock{}, &LoggerInterfaceMock{})

	cuisineType.Name = "deleteError"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/cuisineType/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: cuisineType.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"error"}`), body)
}
