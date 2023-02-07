package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type IngredientHandlersMock struct {
}

// ========= Ingredients =========

func (h *IngredientHandlersMock) GetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}

func (h *IngredientHandlersMock) GetSingle(w http.ResponseWriter, r *http.Request, recipeID string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}

func (h *IngredientHandlersMock) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}

func (h *IngredientHandlersMock) Update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}

func (h *IngredientHandlersMock) Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte(""))
}

// ==================================================================================================

func Test_IngredientNotImplemented(t *testing.T) {
	e := NewIngredientEndpoints(&IngredientHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.NotImplemented(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_IngredientGetAll(t *testing.T) {
	e := NewIngredientEndpoints(&IngredientHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetAll(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_IngredientGetSingle(t *testing.T) {
	e := NewIngredientEndpoints(&IngredientHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetSingle(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_IngredientGetUnits(t *testing.T) {
	e := NewIngredientEndpoints(&IngredientHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetUnits(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_IngredientCreate(t *testing.T) {
	e := NewIngredientEndpoints(&IngredientHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Create(c)

	assert.Equal(t, w.Code, http.StatusCreated)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_IngredientUpdate(t *testing.T) {
	e := NewIngredientEndpoints(&IngredientHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Update(c)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_IngredientDelete(t *testing.T) {
	e := NewIngredientEndpoints(&IngredientHandlersMock{})
	d := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(d)

	e.Delete(c)

	assert.Equal(t, d.Code, http.StatusNoContent)
}
