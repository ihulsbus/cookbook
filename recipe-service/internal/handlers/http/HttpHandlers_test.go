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
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
)

type RecipeServiceMock struct {
}

var (
	recipes []m.RecipeDTO
	recipe  m.RecipeDTO = m.RecipeDTO{
		ID:           uuid.New(),
		Name:         "recipe",
		Description:  "description",
		ServingCount: 1,
	}
)

// ====== RecipeService ======

func (s *RecipeServiceMock) FindAll(pagination m.PaginationRequest) (m.PaginatedResponse[m.RecipeDTO], error) {
	switch recipe.Name {
	case "findall":
		return m.PaginatedResponse[m.RecipeDTO]{
			Data:       recipes,
			Pagination: m.NewPaginationMetadata(pagination.Page, pagination.Limit, int64(len(recipes))),
		}, nil
	default:
		return m.PaginatedResponse[m.RecipeDTO]{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) FindSingle(recipeDTO m.RecipeDTO) (m.RecipeDTO, error) {
	switch recipe.Name {
	case "find":
		return recipe, nil
	case "notfound":
		return m.RecipeDTO{}, errors.New("not found")
	default:
		return m.RecipeDTO{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) Create(recipeDTO m.RecipeDTO) (m.RecipeDTO, error) {
	switch recipeDTO.Name {
	case "create":
		return recipe, nil
	default:
		return m.RecipeDTO{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) Update(recipeDTO m.RecipeDTO) (m.RecipeDTO, error) {
	switch recipeDTO.Name {
	case "update":
		return recipe, nil
	default:
		return m.RecipeDTO{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) Delete(recipeDTO m.RecipeDTO) error {
	switch recipe.Name {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

// Logging
type LoggerInterfaceMock struct{}

func (l *LoggerInterfaceMock) Debugf(format string, args ...interface{}) {}
func (l *LoggerInterfaceMock) Warnf(format string, args ...interface{})  {}
func (l *LoggerInterfaceMock) Errorf(format string, args ...interface{}) {}
func (l *LoggerInterfaceMock) Infof(format string, args ...interface{})  {}

// ==================================================================================================
func TestRecipeGetAll_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recipes = append(recipes, recipe)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "findall"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	pagination := m.NewPaginationMetadata(1, 25, int64(len(recipes)))
	expectedBody, _ := json.Marshal(m.PaginatedResponse[m.RecipeDTO]{Data: recipes, Pagination: pagination})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestRecipeGetAll_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recipes = append(recipes, recipe)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestRecipeGet_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	recipe.Name = "find"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(recipe)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestRecipeGet_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	recipe.Name = "find"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid recipe ID"}`, string(body))
}

func TestRecipeGet_notFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	recipe.Name = "notfound"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"recipe not found"}`, string(body))
}

func TestRecipeGet_FindErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	recipe.Name = "finderr"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestRecipeCreate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	createRecipe := m.RecipeDTO{
		Name: "create",
	}
	reqBody, _ := json.Marshal(createRecipe)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assertBody, _ := json.Marshal(recipe)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, assertBody, body)
}

func TestRecipeCreate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unexpected JSON input"}`, string(body))
}

func TestRecipeCreate_CreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	createRecipe := m.RecipeDTO{
		Name: "error",
	}
	reqBody, _ := json.Marshal(createRecipe)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestRecipeUpdate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "update"
	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestRecipeUpdate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestRecipeUpdate_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid recipe ID"}`, string(body))
}

func TestRecipeUpdate_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "fail"
	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestRecipeDelete_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "delete"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRecipeDelete_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"invalid recipe ID"}`), body)
}

func TestRecipeDelete_DeleteErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "error"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"error"}`), body)
}
