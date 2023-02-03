package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers interface {
	IngredientCreate(w http.ResponseWriter, r *http.Request)
	IngredientDelete(w http.ResponseWriter, r *http.Request)
	IngredientGetAll(w http.ResponseWriter, r *http.Request)
	IngredientGetSingle(w http.ResponseWriter, r *http.Request, ingredientID string)
	RecipeCreate(w http.ResponseWriter, r *http.Request)
	RecipeDelete(w http.ResponseWriter, r *http.Request)
	RecipeGet(w http.ResponseWriter, r *http.Request, recipeID string)
	RecipeGetAll(w http.ResponseWriter, r *http.Request)
	RecipeImageUpload(w http.ResponseWriter, r *http.Request, recipeID string)
	RecipeUpdate(w http.ResponseWriter, r *http.Request)
}

type Endpoints struct {
	handlers Handlers
}

func NewEndpoints(handlers Handlers) *Endpoints {
	return &Endpoints{
		handlers: handlers,
	}
}

func (e Endpoints) RecipeGetAll(ctx *gin.Context) {
	e.handlers.RecipeGetAll(ctx.Writer, ctx.Request)
}

func (e Endpoints) RecipeGet(ctx *gin.Context) {
	e.handlers.RecipeGet(ctx.Writer, ctx.Request, ctx.Param("recipeID"))
}

func (e Endpoints) RecipeUpdate(ctx *gin.Context) {
	e.handlers.RecipeUpdate(ctx.Writer, ctx.Request)
}

func (e Endpoints) RecipeCreate(ctx *gin.Context) {
	e.handlers.RecipeCreate(ctx.Writer, ctx.Request)
}

func (e Endpoints) RecipeDelete(ctx *gin.Context) {
	e.handlers.RecipeDelete(ctx.Writer, ctx.Request)
}

func (e Endpoints) RecipeImageUpload(ctx *gin.Context) {
	// This is dirty, but I do not want gin awareness beyond the endpoints level
	e.handlers.RecipeImageUpload(ctx.Writer, ctx.Request, ctx.Param("recipeID"))
}

func (e Endpoints) IngredientGetAll(ctx *gin.Context) {
	e.handlers.IngredientGetAll(ctx.Writer, ctx.Request)
}

func (e Endpoints) IngredientGetSingle(ctx *gin.Context) {
	e.handlers.IngredientGetSingle(ctx.Writer, ctx.Request, ctx.Param("ingredientID"))
}

func (e Endpoints) IngredientCreate(ctx *gin.Context) {
	e.handlers.IngredientCreate(ctx.Writer, ctx.Request)
}

func (e Endpoints) IngredientDelete(ctx *gin.Context) {
	e.handlers.IngredientDelete(ctx.Writer, ctx.Request)
}
