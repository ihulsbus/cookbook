package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type RecipeHandlersMock struct {
}

func (h *RecipeHandlersMock) GetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *RecipeHandlersMock) Get(w http.ResponseWriter, r *http.Request, recipeID string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *RecipeHandlersMock) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *RecipeHandlersMock) ImageUpload(w http.ResponseWriter, r *http.Request, recipeID string) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *RecipeHandlersMock) Update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *RecipeHandlersMock) Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte(""))
}

// ==================================================================================================

func Test_RecipeNotImplemented(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.NotImplemented(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_RecipeGetAll(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetAll(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_RecipeGet(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Get(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_GetInstruction(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetInstruction(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_RecipeCreate(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Create(c)

	assert.Equal(t, w.Code, http.StatusCreated)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_CreateInstruction(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.CreateInstruction(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_RecipeImageUpload(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.ImageUpload(c)

	assert.Equal(t, w.Code, http.StatusCreated)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}
func Test_RecipeUpdate(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Update(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_UpdateInstruction(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.UpdateInstruction(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_RecipeDelete(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Delete(c)

	assert.Equal(t, w.Code, http.StatusNoContent)
}
