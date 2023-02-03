package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	m "github.com/ihulsbus/cookbook/internal/models"
)

type RecipeServiceMock struct {
}

type ImageServiceMock struct {
}

var (
	recipes []m.Recipe
	recipe  m.Recipe = m.Recipe{
		RecipeName: "create",
	}
)

// ====== RecipeService ======

func (s *RecipeServiceMock) FindAllRecipes() ([]m.Recipe, error) {
	return recipes, nil
}

func (s *RecipeServiceMock) FindSingleRecipe(recipeID int) (m.Recipe, error) {
	if recipeID == 1 {
		return m.Recipe{RecipeName: "recipe1"}, nil
	} else {
		return m.Recipe{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) CreateRecipe(recipe m.Recipe) (m.Recipe, error) {
	if recipe.RecipeName != "create" {
		return recipe, errors.New("error")
	}
	return recipe, nil
}

func (s *RecipeServiceMock) UpdateRecipe(recipe m.Recipe) (m.Recipe, error) {
	if recipe.RecipeName != "update" {
		return recipe, errors.New("error")
	}
	return recipe, nil
}

func (s *RecipeServiceMock) DeleteRecipe(recipe m.Recipe) error {
	if recipe.ID != 1 {
		return errors.New("error")
	}
	return nil
}

// ====== ImageService ======

func (S *ImageServiceMock) UploadImage(file multipart.File, recipeID int) bool {
	if recipeID != 1 {
		return false
	}
	return true
}

// ==================================================================================================

func TestRecipeGetAll_OK(t *testing.T) {
	recipes = append(recipes, m.Recipe{RecipeName: "recipe1"}, m.Recipe{RecipeName: "recipe2"})
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/v1/recipe", nil)
	w := httptest.NewRecorder()

	h.RecipeGetAll(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(recipes)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestRecipeGetAll_Err(t *testing.T) {
	recipes = append(recipes, m.Recipe{RecipeName: "recipe1"}, m.Recipe{RecipeName: "recipe2"})
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/v1/recipe", nil)
	w := httptest.NewRecorder()

	h.RecipeGetAll(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(recipes)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestRecipeGet_OK(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.RecipeGet(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(m.Recipe{RecipeName: "recipe1"})

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestRecipeGet_AtoiErr(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.RecipeGet(w, req, "")

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
}

func TestRecipeGet_FindErr(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.RecipeGet(w, req, "2")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, string(body), `{"code":500,"msg":"Internal Server Error. (error)"}`)
}

func TestRecipeCreate_OK(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("GET", "http://example.com/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.RecipeCreate(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, body, reqBody)
}

func TestRecipeCreate_UnmarshalErr(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/v1/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.RecipeCreate(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (unexpected end of JSON input)"}`))
}

func TestRecipeCreate_CreateErr(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	recipe.RecipeName = "err"
	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("GET", "http://example.com/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.RecipeCreate(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}
