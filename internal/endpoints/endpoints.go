package endpoints

import (
	"github.com/gin-gonic/gin"
	h "github.com/ihulsbus/cookbook/internal/handlers"
)

type Endpoints struct {
	handlers *h.Handlers
}

func NewEndpoints(handlers *h.Handlers) *Endpoints {
	return &Endpoints{
		handlers: handlers,
	}
}

func (e Endpoints) RecipeGetAll(ctx *gin.Context) {
	e.handlers.RecipeGetAll(ctx.Writer, ctx.Request)
}

func (e Endpoints) RecipeGet(ctx *gin.Context) {
	e.handlers.RecipeGetAll(ctx.Writer, ctx.Request)
}

func (e Endpoints) RecipeUpdate(ctx *gin.Context) {
	e.handlers.RecipeGetAll(ctx.Writer, ctx.Request)
}

func (e Endpoints) RecipeCreate(ctx *gin.Context) {
	e.handlers.RecipeGetAll(ctx.Writer, ctx.Request)
}

func (e Endpoints) RecipeDelete(ctx *gin.Context) {
	e.handlers.RecipeGetAll(ctx.Writer, ctx.Request)
}

func (e Endpoints) RecipeImageUpload(ctx *gin.Context) {
	e.handlers.RecipeGetAll(ctx.Writer, ctx.Request)
}

func (e Endpoints) IngredientGetAll(ctx *gin.Context) {
	e.handlers.RecipeGetAll(ctx.Writer, ctx.Request)
}

func (e Endpoints) IngredientGetSingle(ctx *gin.Context) {
	e.handlers.RecipeGetAll(ctx.Writer, ctx.Request)
}

func (e Endpoints) IngredientCreate(ctx *gin.Context) {
	e.handlers.RecipeGetAll(ctx.Writer, ctx.Request)
}

func (e Endpoints) IngredientDelete(ctx *gin.Context) {
	e.handlers.RecipeGetAll(ctx.Writer, ctx.Request)
}
