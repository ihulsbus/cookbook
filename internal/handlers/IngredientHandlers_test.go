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

type IngredientServiceMock struct {
}

var (
	ingredients []m.Ingredient
	ingredient  m.Ingredient = m.Ingredient{
		IngredientName: "ingredient",
	}
)

func (s *IngredientServiceMock) FindAll() ([]m.Ingredient, error) {
	ingredients = append(ingredients, m.Ingredient{IngredientName: "ingredient1"}, m.Ingredient{IngredientName: "ingredient2"})
	return ingredients, nil
}

func (s *IngredientServiceMock) FindSingle(ingredientID int) (m.Ingredient, error) {
	switch ingredientID {
	case 1:
		return ingredient, nil
	default:
		return m.Ingredient{}, errors.New("error")
	}
}

func (s *IngredientServiceMock) Create(ingredient m.Ingredient) (m.Ingredient, error) {
	switch ingredient.IngredientName {
	case "ingredient":
		return ingredient, nil
	default:
		return ingredient, errors.New("error")
	}
}

func (s *IngredientServiceMock) Update(ingredient m.Ingredient) (m.Ingredient, error) {
	switch ingredient.IngredientName {
	case "update":
		return ingredient, nil
	default:
		return ingredient, errors.New("error")
	}
}

func (s *IngredientServiceMock) Delete(ingredient m.Ingredient) error {
	switch ingredient.ID {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

// ==================================================================================================
func TestIngredientGetAll_OK(t *testing.T) {
	var expected []m.Ingredient
	expected = append(expected, m.Ingredient{IngredientName: "ingredient1"}, m.Ingredient{IngredientName: "ingredient2"})
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/v1/ingredients", nil)
	w := httptest.NewRecorder()

	h.IngredientGetAll(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(expected)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestIngredientGet_OK(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/v1/ingredient/1", nil)
	w := httptest.NewRecorder()

	h.IngredientGetSingle(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(ingredient)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestIngredientGet_AtoiErr(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/v1/ingredient/", nil)
	w := httptest.NewRecorder()

	h.IngredientGetSingle(w, req, "")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (strconv.Atoi: parsing \"\": invalid syntax)"}`))
}

func TestIngredientGet_FindErr(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/v1/ingredient/0", nil)
	w := httptest.NewRecorder()

	h.IngredientGetSingle(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestIngredientCreate_OK(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(ingredient)

	req := httptest.NewRequest("POST", "http://example.com/v1/ingredient", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.IngredientCreate(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, body, reqBody)
}

func TestIngredientCreate_UnmarshallErr(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/v1/ingredient", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.IngredientCreate(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (unexpected end of JSON input)"}`))
}

func TestIngredientCreate_CreateErr(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	errIngredient := ingredient
	errIngredient.IngredientName = "err"

	reqBody, _ := json.Marshal(errIngredient)

	req := httptest.NewRequest("POST", "http://example.com/v1/ingredient", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.IngredientCreate(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestIngredientDelete_OK(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	deleteIngredient := ingredient
	deleteIngredient.ID = 1
	deleteIngredient.IngredientName = "ingredient"
	reqBody, _ := json.Marshal(deleteIngredient)

	req := httptest.NewRequest("DELETE", "http://example.com/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.IngredientDelete(w, req)

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)
}

func TestIngredientDelete_UnmarshalErr(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/v1/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.IngredientDelete(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (unexpected end of JSON input)"}`))
}

func TestIngredientDelete_IDRequiredErr(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	deleteIngredient := ingredient
	deleteIngredient.ID = 0
	deleteIngredient.IngredientName = "ingredient"
	reqBody, _ := json.Marshal(deleteIngredient)

	req := httptest.NewRequest("DELETE", "http://example.com/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.IngredientDelete(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestIngredientDelete_DeleteErr(t *testing.T) {
	h := NewHandlers(&RecipeServiceMock{}, &IngredientServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	deleteIngredient := ingredient
	deleteIngredient.ID = 2
	deleteIngredient.IngredientName = "ingredient"
	reqBody, _ := json.Marshal(deleteIngredient)

	req := httptest.NewRequest("DELETE", "http://example.com/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.IngredientDelete(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}
