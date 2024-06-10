package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"

	m "recipe-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type RecipeHandlersMock struct {
}

type MiddlewareMock struct{}

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
func (h *RecipeHandlersMock) Create(user *m.User, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *RecipeHandlersMock) Update(w http.ResponseWriter, r *http.Request, recipeID string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *RecipeHandlersMock) Delete(w http.ResponseWriter, r *http.Request, recipeID string) {
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte(""))
}

func (MiddlewareMock) UserFromContext(ctx *gin.Context) (*m.User, error) {
	return &m.User{}, nil
}

// ==================================================================================================
func Test_RecipeGetAll(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{}, &MiddlewareMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetAll(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_RecipeGet(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{}, &MiddlewareMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Get(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_RecipeCreate(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{}, &MiddlewareMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Create(c)

	assert.Equal(t, w.Code, http.StatusCreated)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_RecipeUpdate(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{}, &MiddlewareMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Update(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_RecipeDelete(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{}, &MiddlewareMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Delete(c)

	assert.Equal(t, w.Code, http.StatusNoContent)
}
