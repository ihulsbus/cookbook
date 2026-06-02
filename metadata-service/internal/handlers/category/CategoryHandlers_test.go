package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
)

type CategoryServiceMock struct {
}

var (
	categories []models.CategoryDTO
	category   models.CategoryDTO = models.CategoryDTO{
		ID:   uuid.New(),
		Name: "category",
	}
)

// ====== CategoryService ======

func (s *CategoryServiceMock) FindAll(pagination models.PaginationRequest) (models.PaginatedResponse[models.CategoryDTO], error) {
	switch category.Name {
	case "findall":
		return models.PaginatedResponse[models.CategoryDTO]{
			Data:       categories,
			Pagination: models.NewPaginationMetadata(pagination.Page, pagination.Limit, int64(len(categories))),
		}, nil
	default:
		return models.PaginatedResponse[models.CategoryDTO]{}, errors.New("error")
	}
}

func (s *CategoryServiceMock) FindSingle(categoryDTO models.CategoryDTO) (models.CategoryDTO, error) {
	switch category.Name {
	case "find":
		return category, nil
	case "notfound":
		return models.CategoryDTO{}, errors.New("not found")
	default:
		return models.CategoryDTO{}, errors.New("error")
	}
}

func (s *CategoryServiceMock) Create(categoryDTO models.CategoryDTO) (models.CategoryDTO, error) {
	switch categoryDTO.Name {
	case "create":
		return category, nil
	default:
		return models.CategoryDTO{}, errors.New("error")
	}
}

func (s *CategoryServiceMock) Update(categoryDTO models.CategoryDTO) (models.CategoryDTO, error) {
	switch categoryDTO.Name {
	case "update":
		return category, nil
	default:
		return models.CategoryDTO{}, errors.New("error")
	}
}

func (s *CategoryServiceMock) Delete(categoryDTO models.CategoryDTO) error {
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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

	category.Name = "findall"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/category", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	pagination := models.NewPaginationMetadata(1, 25, int64(len(categories)))
	expectedResponse := models.PaginatedResponse[models.CategoryDTO]{Data: categories, Pagination: pagination}
	expectedBody, _ := json.Marshal(expectedResponse)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestCategoryGetAll_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	categories = append(categories, category)
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

	createCategory := models.CategoryDTO{
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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

	createCategory := models.CategoryDTO{
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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewCategoryHandlers(&CategoryServiceMock{}, &models.LoggerInterfaceMock{})

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
