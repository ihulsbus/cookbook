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

type CategoryServiceMock struct {
}

var (
	categories []m.CategoryDTO
	category   m.CategoryDTO = m.CategoryDTO{
		ID:   uuid.New(),
		Name: "category",
	}
)

// ====== CategoryService ======

func (s *CategoryServiceMock) FindAll() ([]m.CategoryDTO, error) {
	switch category.Name {
	case "findall":
		return categories, nil
	case "notfound":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (s *CategoryServiceMock) FindSingle(categoryDTO m.CategoryDTO) (m.CategoryDTO, error) {
	switch category.Name {
	case "find":
		return category, nil
	case "notfound":
		return m.CategoryDTO{}, errors.New("not found")
	default:
		return m.CategoryDTO{}, errors.New("error")
	}
}

func (s *CategoryServiceMock) Create(categoryDTO m.CategoryDTO) (m.CategoryDTO, error) {
	switch categoryDTO.Name {
	case "create":
		return category, nil
	default:
		return m.CategoryDTO{}, errors.New("error")
	}
}

func (s *CategoryServiceMock) Update(categoryDTO m.CategoryDTO) (m.CategoryDTO, error) {
	switch categoryDTO.Name {
	case "update":
		return category, nil
	default:
		return m.CategoryDTO{}, errors.New("error")
	}
}

func (s *CategoryServiceMock) Delete(categoryDTO m.CategoryDTO) error {
	switch category.Name {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

func TestCategoryGetAll_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	categories = append(categories, category)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	category.Name = "findall"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/category", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(categories)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestCategoryGetAll_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	categories = append(categories, category)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	category.Name = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/category", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no categories found"}`, string(body))
}

func TestCategoryGetAll_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	categories = append(categories, category)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	category.Name = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/category", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestCategoryGet_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	category.Name = "find"
	req := httptest.NewRequest("GET", "http://example.com/api/v2/category/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: category.ID.String()},
	}

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(category)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestCategoryGet_ID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	category.Name = "finderr"
	req := httptest.NewRequest("GET", "http://example.com/api/v2/category/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid category ID"}`, string(body))
}

func TestCategoryGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	category.Name = "notfound"
	req := httptest.NewRequest("GET", "http://example.com/api/v2/category/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: category.ID.String()},
	}

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"category not found"}`, string(body))
}

func TestCategoryGet_FindErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	category.Name = "finderr"
	req := httptest.NewRequest("GET", "http://example.com/api/v2/category/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: category.ID.String()},
	}

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestCategoryCreate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	createCategory := m.CategoryDTO{
		Name: "create",
	}
	reqBody, _ := json.Marshal(createCategory)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/category/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assertBody, _ := json.Marshal(category)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, assertBody, body)
}

func TestCategoryCreate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/category/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestCategoryCreate_CreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	createCategory := m.CategoryDTO{
		Name: "createError",
	}
	reqBody, _ := json.Marshal(createCategory)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/category/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestCategoryUpdate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	category.Name = "update"
	reqBody, _ := json.Marshal(category)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/category/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: category.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestCategoryUpdate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/category/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: category.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestCategoryUpdate_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(category)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/category/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid category ID"}`, string(body))
}

func TestCategoryUpdate_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	category.Name = "updateFail"
	reqBody, _ := json.Marshal(category)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/category/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: category.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestCategoryDelete_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	category.Name = "delete"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/category/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: category.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, ``, string(body))
}

func TestCategoryDelete_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/category/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid category ID"}`, string(body))
}

func TestCategoryDelete_DeleteErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	category.Name = "deleteError"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/category/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: category.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}
