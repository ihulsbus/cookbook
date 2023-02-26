package endpoints

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	m "github.com/ihulsbus/cookbook/internal/models"
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
func (h *RecipeHandlersMock) GetIngredients(w http.ResponseWriter, r *http.Request, recipeID string) {
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
func (h *RecipeHandlersMock) ImageUpload(w http.ResponseWriter, r *http.Request, recipeID string) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *RecipeHandlersMock) GetInstruction(w http.ResponseWriter, r *http.Request, recipeID string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *RecipeHandlersMock) CreateInstruction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *RecipeHandlersMock) UpdateInstruction(w http.ResponseWriter, r *http.Request, recipeID string) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *RecipeHandlersMock) DeleteInstruction(w http.ResponseWriter, r *http.Request, recipeID string) {
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte(""))
}

func (MiddlewareMock) UserFromContext(ctx context.Context) (*m.User, error) {
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

func Test_RecipeGetIngredients(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{}, &MiddlewareMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetIngredients(c)

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

func Test_RecipeImageUpload(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{}, &MiddlewareMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.ImageUpload(c)

	assert.Equal(t, w.Code, http.StatusCreated)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_GetInstruction(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{}, &MiddlewareMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetInstruction(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_CreateInstruction(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{}, &MiddlewareMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.CreateInstruction(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_UpdateInstruction(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{}, &MiddlewareMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.UpdateInstruction(c)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_DeleteInstruction(t *testing.T) {
	e := NewRecipeEndpoints(&RecipeHandlersMock{}, &MiddlewareMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.DeleteInstruction(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
