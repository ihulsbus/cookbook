package endpoints

import (
	"net/http"

	m "recipe-service/internal/models"

	"github.com/gin-gonic/gin"
)

type RecipeHandlers interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, recipeID string)
	Create(user *m.User, w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request, recipeID string)
	Delete(w http.ResponseWriter, r *http.Request, recipeID string)
}

type Middleware interface {
	UserFromContext(ctx *gin.Context) (*m.User, error)
}

type RecipeEndpoints struct {
	handlers   RecipeHandlers
	middleware Middleware
}

func NewRecipeEndpoints(handlers RecipeHandlers, middleware Middleware) *RecipeEndpoints {
	return &RecipeEndpoints{
		handlers:   handlers,
		middleware: middleware,
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
// @Param			id	path		int	true	"Recipe ID"
// @Success		200	{object}	models.Recipe
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id} [get]
func (e RecipeEndpoints) Get(ctx *gin.Context) {
	e.handlers.Get(ctx.Writer, ctx.Request, ctx.Param("id"))
}

// @Summary		Create a recipe
// @Description	Creates a new recipe and returns a JSON object of the created recipe
// @tags			recipes
// @Accept		json
// @Produce		json
// @Success		200	{object}	models.Recipe
// @Param			requestbody	body		models.Recipe	true	"Create a recipe"
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe [post]
func (e RecipeEndpoints) Create(ctx *gin.Context) {
	user, err := e.middleware.UserFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	e.handlers.Create(user, ctx.Writer, ctx.Request)
}

// @Summary		Update a recipe
// @Description	Updates a single recipe and return the JSON object of the updated recipe
// @Tags			recipes
// @Accept		json
// @Produce		json
// @Param			id	path		int	true	"Recipe ID"
// @Param			requestbody	body		models.Recipe	true	"Update a recipe"
// @Success		200	{object}	models.Recipe
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id} [put]
func (e RecipeEndpoints) Update(ctx *gin.Context) {
	e.handlers.Update(ctx.Writer, ctx.Request, ctx.Param("id"))
}

// @Summary		Delete a recipe
// @Description	Deletes a recipe. Returns a simple http status code
// @Tags			recipes
// @Produce		json
// @Param			id	path	int	true	"Recipe ID"
// @Success		204
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id} [delete]
func (e RecipeEndpoints) Delete(ctx *gin.Context) {
	e.handlers.Delete(ctx.Writer, ctx.Request, ctx.Param("id"))
}
