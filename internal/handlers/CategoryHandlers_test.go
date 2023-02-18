package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	m "github.com/ihulsbus/cookbook/internal/models"
	"github.com/stretchr/testify/assert"
)

type CategoryServiceMock struct {
}

var (
	categorys []m.Category
	category  m.Category = m.Category{
		CategoryName: "category",
	}
)

// ====== CategoryService ======

func (s *CategoryServiceMock) FindAll() ([]m.Category, error) {
	return categorys, nil
}

func (s *CategoryServiceMock) FindSingle(categoryID uint) (m.Category, error) {
	switch categoryID {
	case 1:
		return m.Category{CategoryName: "category1"}, nil
	case 2:
		return m.Category{CategoryName: "category2"}, nil
	default:
		return m.Category{}, errors.New("error")
	}
}

func (s *CategoryServiceMock) Create(category m.Category) (m.Category, error) {
	switch category.CategoryName {
	case "category":
		return category, nil
	default:
		return category, errors.New("error")
	}
}

func (s *CategoryServiceMock) Update(category m.Category, categoryID uint) (m.Category, error) {
	switch category.CategoryName {
	case "category":
		return category, nil
	default:
		return category, errors.New("error")
	}
}

func (s *CategoryServiceMock) Delete(categoryID uint) error {
	switch categoryID {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

func TestCategoryGetAll_OK(t *testing.T) {
	categorys = append(categorys, m.Category{CategoryName: "category1"}, m.Category{CategoryName: "category2"})
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/category", nil)
	w := httptest.NewRecorder()

	h.GetAll(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(categorys)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestCategoryGet_OK(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/category/1", nil)
	w := httptest.NewRecorder()

	h.Get(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(m.Category{CategoryName: "category1"})

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestCategoryGet_AtoiErr(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/category/1", nil)
	w := httptest.NewRecorder()

	h.Get(w, req, "")

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
}

func TestCategoryGet_FindErr(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/category/1", nil)
	w := httptest.NewRecorder()

	h.Get(w, req, "3")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, `{"code":500,"msg":"Internal Server Error. (error)"}`, string(body))
}

func TestCategoryCreate_OK(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(category)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/category/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, body, reqBody)
}

func TestCategoryCreate_UnmarshalErr(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v1/category/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (unexpected end of JSON input)"}`))
}

func TestCategoryCreate_CreateErr(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	category.CategoryName = "err"
	reqBody, _ := json.Marshal(category)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/category/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestCategoryUpdate_OK(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	updateCategory := category
	updateCategory.ID = 1
	updateCategory.CategoryName = "category"
	reqBody, _ := json.Marshal(updateCategory)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/category/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, reqBody)
}

func TestCategoryUpdate_AtoiErr(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/category/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Update(w, req, "")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, []byte(`{"code":500,"msg":"Internal server error"}`), body)
}

func TestCategoryUpdate_UnmarshalErr(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/category/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (unexpected end of JSON input)"}`))
}

func TestCategoryUpdate_IDRequiredErr(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	updateCategory := category
	updateCategory.ID = 0
	updateCategory.CategoryName = "category"
	reqBody, _ := json.Marshal(updateCategory)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/category/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestCategoryUpdate_UpdateErr(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	updateCategory := category
	updateCategory.ID = 2
	updateCategory.CategoryName = "fail"
	reqBody, _ := json.Marshal(updateCategory)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/category/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestCategoryDelete_OK(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/category/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "1")

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)
}

func TestCategoryDelete_AtoiErr(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/category/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "")

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestCategoryDelete_IDRequiredErr(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/category/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestCategoryDelete_DeleteErr(t *testing.T) {
	h := NewCategoryHandlers(&CategoryServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/category/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "2")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`), body)
}
