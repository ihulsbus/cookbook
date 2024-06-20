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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type IngredientServiceMock struct {
}

var (
	ingredients []m.IngredientDTO
	ingredient  m.IngredientDTO = m.IngredientDTO{
		ID:   uuid.New(),
		Name: "ingredient",
	}

	units []m.UnitDTO
	unit  m.UnitDTO = m.UnitDTO{
		ID:        uuid.New(),
		FullName:  "Fluid Ounce",
		ShortName: "fl oz",
	}
)

func (s *IngredientServiceMock) FindAll() ([]m.IngredientDTO, error) {
	switch ingredient.Name {
	case "find":
		ingredients = append(ingredients, ingredient)
		return ingredients, nil
	case "notfound":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (s *IngredientServiceMock) FindUnits() ([]m.UnitDTO, error) {
	units = append(units, unit)
	return units, nil
}

func (s *IngredientServiceMock) FindSingle(ingredientDTO m.IngredientDTO) (m.IngredientDTO, error) {
	switch ingredientDTO.Name {
	case "find":
		return ingredient, nil
	case "notfound":
		return m.IngredientDTO{}, errors.New("not found")
	default:
		return m.IngredientDTO{}, errors.New("error")
	}
}

func (s *IngredientServiceMock) Create(ingredientDTO m.IngredientDTO) (m.IngredientDTO, error) {
	switch ingredient.Name {
	case "create":
		return ingredient, nil
	default:
		return ingredient, errors.New("error")
	}
}

func (s *IngredientServiceMock) Update(ingredientDTO m.IngredientDTO) (m.IngredientDTO, error) {
	switch ingredient.Name {
	case "update":
		return ingredient, nil
	default:
		return ingredient, errors.New("error")
	}
}

func (s *IngredientServiceMock) Delete(ingredientDTO m.IngredientDTO) error {
	switch ingredientDTO.Name {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

// ==================================================================================================
func TestIngredientGetAll_OK(t *testing.T) {
	ingredients = append(ingredients, ingredient)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	ingredient.Name = "findall"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredients", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(ingredients)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestIngredientGetUnits_OK(t *testing.T) {
	units = append(units, unit)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredient/units", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetUnits(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(units)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestIngredientGet_OK(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredient/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetSingle(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(ingredient)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestIngredientGet_AtoiErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredient/", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetSingle(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (strconv.Atoi: parsing \"\": invalid syntax)"}`))
}

func TestIngredientGet_FindErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredient/0", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetSingle(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestIngredientCreate_OK(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(ingredient)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/ingredient", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, body, reqBody)
}

func TestIngredientCreate_UnmarshallErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/ingredient", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (unexpected end of JSON input)"}`))
}

func TestIngredientCreate_CreateErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	errIngredient := ingredient
	errIngredient.Name = "err"

	reqBody, _ := json.Marshal(errIngredient)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/ingredient", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestIngredientUpdate_OK(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(ingredient)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/ingredient/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestIngredientUpdate_AtoiErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredient/", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, []byte(`{"code":500,"msg":"Internal server error"}`), body)
}

func TestIngredientUpdate_UnmarshalErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/ingredient/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (unexpected end of JSON input)"}`))
}

func TestIngredientUpdate_IDRequiredErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(ingredient)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/ingredient/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestIngredientUpdate_UpdateErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	ingredient.Name = "fail"
	reqBody, _ := json.Marshal(ingredient)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/ingredient/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestIngredientDelete_OK(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/ingredient/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)
}

func TestIngredientDelete_AtoiErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredient/", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, []byte(`{"code":500,"msg":"Internal server error"}`), body)
}

func TestIngredientDelete_IDRequiredErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/ingredient/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestIngredientDelete_DeleteErr(t *testing.T) {
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/ingredient/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`), body)
}
