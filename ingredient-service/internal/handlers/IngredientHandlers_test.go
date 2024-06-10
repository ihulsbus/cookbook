package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	m "ingredient-service/internal/models"

	"github.com/stretchr/testify/assert"
)

type IngredientServiceMock struct {
}

var (
	ingredients []m.Ingredient
	units       []m.Unit
	ingredient  m.Ingredient = m.Ingredient{
		IngredientName: "ingredient",
	}
	unit m.Unit = m.Unit{
		ID:        1,
		FullName:  "Fluid Ounce",
		ShortName: "fl oz"}
)

func (s *IngredientServiceMock) FindAll() ([]m.Ingredient, error) {
	ingredients = append(ingredients, m.Ingredient{IngredientName: "ingredient1"}, m.Ingredient{IngredientName: "ingredient2"})
	return ingredients, nil
}

func (s *IngredientServiceMock) FindUnits() ([]m.Unit, error) {
	units = append(units, unit, m.Unit{ID: 2, FullName: "Ounce", ShortName: "oz"})
	return units, nil
}

func (s *IngredientServiceMock) FindSingle(ingredientID uint) (m.Ingredient, error) {
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

func (s *IngredientServiceMock) Update(ingredient m.Ingredient, ingredientID uint) (m.Ingredient, error) {
	switch ingredient.IngredientName {
	case "update":
		return ingredient, nil
	default:
		return ingredient, errors.New("error")
	}
}

func (s *IngredientServiceMock) Delete(ingredientID uint) error {
	switch ingredientID {
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
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/ingredients", nil)
	w := httptest.NewRecorder()

	h.GetAll(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(expected)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestIngredientGetUnits_OK(t *testing.T) {
	var expected []m.Unit
	expected = append(expected, unit, m.Unit{ID: 2, FullName: "Ounce", ShortName: "oz"})
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/ingredient/units", nil)
	w := httptest.NewRecorder()

	h.GetUnits(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(expected)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestIngredientGet_OK(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/ingredient/1", nil)
	w := httptest.NewRecorder()

	h.GetSingle(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(ingredient)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestIngredientGet_AtoiErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/ingredient/", nil)
	w := httptest.NewRecorder()

	h.GetSingle(w, req, "")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (strconv.Atoi: parsing \"\": invalid syntax)"}`))
}

func TestIngredientGet_FindErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/ingredient/0", nil)
	w := httptest.NewRecorder()

	h.GetSingle(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestIngredientCreate_OK(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(ingredient)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/ingredient", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, body, reqBody)
}

func TestIngredientCreate_UnmarshallErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v1/ingredient", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (unexpected end of JSON input)"}`))
}

func TestIngredientCreate_CreateErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	errIngredient := ingredient
	errIngredient.IngredientName = "err"

	reqBody, _ := json.Marshal(errIngredient)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/ingredient", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestIngredientUpdate_OK(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	updateIngredient := ingredient
	updateIngredient.ID = 1
	updateIngredient.IngredientName = "update"
	reqBody, _ := json.Marshal(updateIngredient)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/ingredient/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestIngredientUpdate_AtoiErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/ingredient/", nil)
	w := httptest.NewRecorder()

	h.Update(w, req, "")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, []byte(`{"code":500,"msg":"Internal server error"}`), body)
}

func TestIngredientUpdate_UnmarshalErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/ingredient/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (unexpected end of JSON input)"}`))
}

func TestIngredientUpdate_IDRequiredErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	updateIngredient := ingredient
	updateIngredient.ID = 0
	updateIngredient.IngredientName = "ingredient"
	reqBody, _ := json.Marshal(updateIngredient)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/ingredient/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestIngredientUpdate_UpdateErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	updateIngredient := ingredient
	updateIngredient.ID = 2
	updateIngredient.IngredientName = "fail"
	reqBody, _ := json.Marshal(updateIngredient)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/ingredient/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "2")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestIngredientDelete_OK(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/ingredient/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "1")

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)
}

func TestIngredientDelete_AtoiErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/ingredient/", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, []byte(`{"code":500,"msg":"Internal server error"}`), body)
}

func TestIngredientDelete_IDRequiredErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/ingredient/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestIngredientDelete_DeleteErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/ingredient/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "2")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`), body)
}
