package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type HandlersMock struct {
	// ingredientCreate    func(w http.ResponseWriter, r *http.Request)
	// ingredientDelete    func(w http.ResponseWriter, r *http.Request)
	// ingredientGetAll    func(w http.ResponseWriter, r *http.Request)
	// ingredientGetSingle func(w http.ResponseWriter, r *http.Request, ingredientID string)
	// notImplemented      func(w http.ResponseWriter, r *http.Request)
	// recipeCreate        func(w http.ResponseWriter, r *http.Request)
	// recipeDelete        func(w http.ResponseWriter, r *http.Request)
	// recipeGet           func(w http.ResponseWriter, r *http.Request, recipeID string)
	// recipeGetAll        func(w http.ResponseWriter, r *http.Request)
	// recipeImageUpload   func(w http.ResponseWriter, r *http.Request, recipeID string)
	// recipeUpdate        func(w http.ResponseWriter, r *http.Request)
}

func (h *HandlersMock) IngredientCreate(w http.ResponseWriter, r *http.Request) { return }
func (h *HandlersMock) IngredientDelete(w http.ResponseWriter, r *http.Request) { return }
func (h *HandlersMock) IngredientGetAll(w http.ResponseWriter, r *http.Request) { return }
func (h *HandlersMock) IngredientGetSingle(w http.ResponseWriter, r *http.Request, recipeID string) {
	return
}
func (h *HandlersMock) NotImplemented(w http.ResponseWriter, r *http.Request)             { return }
func (h *HandlersMock) RecipeCreate(w http.ResponseWriter, r *http.Request)               { return }
func (h *HandlersMock) RecipeDelete(w http.ResponseWriter, r *http.Request)               { return }
func (h *HandlersMock) RecipeGet(w http.ResponseWriter, r *http.Request, recipeID string) { return }
func (h *HandlersMock) RecipeGetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{}"))
}
func (h *HandlersMock) RecipeImageUpload(w http.ResponseWriter, r *http.Request, recipeID string) {
	return
}
func (h *HandlersMock) RecipeUpdate(w http.ResponseWriter, r *http.Request) { return }

func Test_RecipeGetAll(t *testing.T) {
	e := NewEndpoints(&HandlersMock{})
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	e.RecipeGetAll(c)

	// assert.Equal(t,)

}
