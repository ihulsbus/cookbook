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

func (e RecipeEndpoints) NotImplemented(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(501, "not implemented")
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
	// This is dirty, but I do not want gin awareness beyond the endpoints level
	e.handlers.Get(ctx.Writer, ctx.Request, ctx.Param("id"))
}

// @Summary		Get a recipe's instruction text
// @Description	Returns the JSON object of the recipe's instructions
// @tags			recipes
// @Produce		json
// @Param			id			path		int				true	"Recipe ID"
// @Success		200			{object}	models.Instruction
// @Failure		401			{string}	string	"unauthorized"
// @Failure		404			{string}	string	"not found"
// @Failure		500			{string}	string	"Any error"
// @Router			/recipe/{id}/instruction [get]
func (e RecipeEndpoints) GetInstruction(ctx *gin.Context) {
	e.NotImplemented(ctx)
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
	e.handlers.Create(ctx.Writer, ctx.Request)
}

// @Summary		Create a recipe's instruction text
// @Description	Creates the instruction text for a recipe and returns the JSON object of the created instructions
// @tags			recipes
// @Accept		json
// @Produce		json
// @Param			id	path		int	true	"Recipe ID"
// @Param			requestbody	body		models.Instruction	true	"Create an instruction"
// @Success		200	{object}	models.Instruction
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id}/instruction [post]
func (e RecipeEndpoints) CreateInstruction(ctx *gin.Context) {
	e.NotImplemented(ctx)
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
	e.handlers.Update(ctx.Writer, ctx.Request)
}

// @Summary		Update a recipe's instruction text
// @Description	Updates the instruction text for a recipe and returns the JSON object of the updated instructions
// @tags			recipes
// @Accept		json
// @Produce		json
// @Param			id	path		int	true	"Recipe ID"
// @Param			requestbody	body		models.Instruction	true	"Update an instruction"
// @Success		200	{object}	models.Instruction
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id}/instruction [put]
func (e RecipeEndpoints) UpdateInstruction(ctx *gin.Context) {
	e.NotImplemented(ctx)
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
	e.handlers.Delete(ctx.Writer, ctx.Request)
}

// @Summary		Upload a recipe image
// @Description	Upload a recipe image used in the frontend. Returns a simple http status code
// @Tags			recipes
// @Produce		json
// @Accept		image/jpeg
// @Param			id	path	int	true	"Recipe ID"
// @Success		201
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id}/cover [post]
func (e RecipeEndpoints) ImageUpload(ctx *gin.Context) {
	// This is dirty, but I do not want gin awareness beyond the endpoints level
	e.handlers.ImageUpload(ctx.Writer, ctx.Request, ctx.Param("id"))
}
