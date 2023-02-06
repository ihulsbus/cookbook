package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecipeHandlers interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, recipeID string)
	Create(w http.ResponseWriter, r *http.Request)
	ImageUpload(w http.ResponseWriter, r *http.Request, recipeID string)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type RecipeEndpoints struct {
	handlers RecipeHandlers
}

func NewRecipeEndpoints(handlers RecipeHandlers) *RecipeEndpoints {
	return &RecipeEndpoints{
		handlers: handlers,
	}
}

// @Summary		Get a list of all available recipes
// @Description	Returns a JSON array of all available recipes
// @tags			recipes
// @produce		json
// @Success		200	{array}		models.Recipe
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe [get]
func (e RecipeEndpoints) GetAll(ctx *gin.Context) {
	e.handlers.GetAll(ctx.Writer, ctx.Request)
}

// @Summary		Get a single recipes
// @Description	Returns a JSON object of a single recipe
// @tags			recipes
// @produce		json
// @Param			id	path	int true "Recipe ID"
// @Success		200	{object}	models.Recipe
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id} [get]
func (e RecipeEndpoints) Get(ctx *gin.Context) {
	// This is dirty, but I do not want gin awareness beyond the endpoints level
	e.handlers.Get(ctx.Writer, ctx.Request, ctx.Param("recipeID"))
}

// @Summary		Create a recipe
// @Description	Creates a new recipe and returns a JSON object of the created recipe
// @tags			recipes
// @Accept		json
// @Produce		json
// @Success		200	{object}	models.Recipe
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe [post]
func (e RecipeEndpoints) Create(ctx *gin.Context) {
	e.handlers.Create(ctx.Writer, ctx.Request)
}

// @Summary		Update a recipe
// @Description	Updates a single recipe and return the JSON object of the updated recipe
// @Tags			recipes
// @Accept		json
// @Produce		json
// @Param			id	path	int true "Recipe ID"
// @Success		200	{object} models.Recipe
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id} [put]
func (e RecipeEndpoints) Update(ctx *gin.Context) {
	e.handlers.Update(ctx.Writer, ctx.Request)
}

// @Summary		Delete a recipe
// @Description	Deletes a recipe. Returns a simple http status code
// @Tags			recipes
// @Produce		json
// @Param			id	path	int true "Recipe ID"
// @Success		204
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id} [delete]
func (e RecipeEndpoints) Delete(ctx *gin.Context) {
	e.handlers.Delete(ctx.Writer, ctx.Request)
}

// @Summary		Upload a recipe image
// @Description	Upload a recipe image used in the frontend. Returns a simple http status code
// @Tags			recipes
// @Produce		json
// @Accept		image/jpeg
// @Param			id	path	int true "Recipe ID"
// @Success		201
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id}/cover [post]
func (e RecipeEndpoints) ImageUpload(ctx *gin.Context) {
	// This is dirty, but I do not want gin awareness beyond the endpoints level
	e.handlers.ImageUpload(ctx.Writer, ctx.Request, ctx.Param("recipeID"))
}
