package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	m "recipe-service/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
)

type RecipeServiceMock struct {
}

var (
	recipes []m.Recipe
	recipe  m.Recipe = m.Recipe{
		RecipeName: "recipe",
	}
)

// ====== RecipeService ======

func (s *RecipeServiceMock) FindAll() ([]m.Recipe, error) {
	return recipes, nil
}

func (s *RecipeServiceMock) FindSingle(recipeID uint) (m.Recipe, error) {
	switch recipeID {
	case 1:
		return m.Recipe{RecipeName: "recipe1"}, nil
	case 2:
		return m.Recipe{RecipeName: "recipe2"}, nil
	default:
		return m.Recipe{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) Create(recipe m.Recipe) (m.Recipe, error) {
	switch recipe.RecipeName {
	case "recipe":
		return recipe, nil
	default:
		return recipe, errors.New("error")
	}
}

func (s *RecipeServiceMock) Update(recipe m.Recipe, recipeID uint) (m.Recipe, error) {
	switch recipe.RecipeName {
	case "recipe":
		return recipe, nil
	default:
		return recipe, errors.New("error")
	}
}

func (s *RecipeServiceMock) Delete(recipeID uint) error {
	switch recipeID {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

// ==================================================================================================

func TestRecipeGetAll_OK(t *testing.T) {
	recipes = append(recipes, m.Recipe{RecipeName: "recipe1"}, m.Recipe{RecipeName: "recipe2"})
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe", nil)
	w := httptest.NewRecorder()

	h.GetAll(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(recipes)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestRecipeGet_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Get(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(m.Recipe{RecipeName: "recipe1"})

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestRecipeGet_AtoiErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Get(w, req, "")

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
}

func TestRecipeGet_FindErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Get(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, `{"code":500,"msg":"Internal Server Error. (error)"}`, string(body))
}

func TestRecipeCreate_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Create(&ginkeycloak.KeyCloakToken{}, w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, body, reqBody)
}

func TestRecipeCreate_UnmarshalErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Create(&ginkeycloak.KeyCloakToken{}, w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (unexpected end of JSON input)"}`))
}

func TestRecipeCreate_CreateErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.RecipeName = "err"
	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Create(&ginkeycloak.KeyCloakToken{}, w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestRecipeUpdate_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	updateRecipe := recipe
	updateRecipe.ID = 1
	updateRecipe.RecipeName = "recipe"
	reqBody, _ := json.Marshal(updateRecipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, reqBody)
}

func TestRecipeUpdate_UnmarshalErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (unexpected end of JSON input)"}`))
}

func TestRecipeUpdate_IDRequiredErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	updateRecipe := recipe
	updateRecipe.ID = 0
	updateRecipe.RecipeName = "recipe"
	reqBody, _ := json.Marshal(updateRecipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestRecipeUpdate_UpdateErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	updateRecipe := recipe
	updateRecipe.ID = 2
	updateRecipe.RecipeName = "fail"
	reqBody, _ := json.Marshal(updateRecipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestRecipeDelete_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "1")

	resp := w.Result()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestRecipeDelete_AtoiErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "")

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestRecipeDelete_IDRequiredErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestRecipeDelete_DeleteErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "2")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}
