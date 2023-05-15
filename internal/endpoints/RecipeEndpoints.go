package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	m "github.com/ihulsbus/cookbook/internal/models"
)

type RecipeHandlers interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, recipeID string)
	Create(user *m.User, w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request, recipeID string)
	Delete(w http.ResponseWriter, r *http.Request, recipeID string)

	ImageUpload(w http.ResponseWriter, r *http.Request, recipeID string)

	GetInstruction(w http.ResponseWriter, r *http.Request, recipeID string)
	CreateInstruction(w http.ResponseWriter, r *http.Request, recipeID string)
	UpdateInstruction(w http.ResponseWriter, r *http.Request, recipeID string)
	DeleteInstruction(w http.ResponseWriter, r *http.Request, recipeID string)

	GetIngredientLink(w http.ResponseWriter, r *http.Request, recipeID string)
	CreateIngredientLink(w http.ResponseWriter, r *http.Request, recipeID string)
	UpdateIngredientLink(w http.ResponseWriter, r *http.Request, recipeID string)
	DeleteIngredientLink(w http.ResponseWriter, r *http.Request, recipeID string)
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
	e.handlers.GetInstruction(ctx.Writer, ctx.Request, ctx.Param("id"))
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
	e.handlers.CreateInstruction(ctx.Writer, ctx.Request, ctx.Param("id"))
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
	e.handlers.UpdateInstruction(ctx.Writer, ctx.Request, ctx.Param("id"))
}

// @Summary		delete a recipe's instruction text
// @Description	Updates the instruction text for a recipe and returns the JSON object of the updated instructions
// @tags			recipes
// @Accept		json
// @Produce		json
// @Param			id	path		int	true	"Recipe ID"
// @Success		204
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id}/instruction [delete]
func (e RecipeEndpoints) DeleteInstruction(ctx *gin.Context) {
	e.handlers.DeleteInstruction(ctx.Writer, ctx.Request, ctx.Param("id"))
}

// @Summary		Get a recipe's ingredients
// @Description	Returns a JSON object with the ingredients and details belonging to a recipe
// @tags			recipes
// @produce		json
// @Param			id	path		int	true	"Recipe ID"
// @Success		200	{array}		models.RecipeIngredient
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id}/ingredients [get]
func (e RecipeEndpoints) GetIngredientLink(ctx *gin.Context) {
	e.handlers.GetIngredientLink(ctx.Writer, ctx.Request, ctx.Param("id"))
}

// @Summary		Create a recipe's ingredient links
// @Description	Creates the ingredient links for a recipe and returns the JSON object of the created ingredient links
// @tags			recipes
// @Accept		json
// @Produce		json
// @Param			id	path		int	true	"Recipe ID"
// @Param			requestbody	body		models.RecipeIngredient	true	"Create an ingredient link"
// @Success		200	{array}		models.Instruction
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id}/ingredients [post]
func (e RecipeEndpoints) CreateIngredientLink(ctx *gin.Context) {
	e.handlers.CreateIngredientLink(ctx.Writer, ctx.Request, ctx.Param("id"))
}

// @Summary		Update a recipe's ingredient links
// @Description	Updates the ingredient links for a recipe and returns the JSON object of the updated ingredient links
// @tags			recipes
// @Accept		json
// @Produce		json
// @Param			id	path		int	true	"Recipe ID"
// @Param			requestbody	body		models.RecipeIngredient	true	"Update an ingredient"
// @Success		200	{array}		models.Instruction
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id}/ingredients [put]
func (e RecipeEndpoints) UpdateIngredientLink(ctx *gin.Context) {
	e.handlers.UpdateIngredientLink(ctx.Writer, ctx.Request, ctx.Param("id"))
}

// @Summary		delete a recipe's ingredient links
// @Description	Updates the ingredient link for a recipe and returns the JSON object of the updated ingredient links
// @tags			recipes
// @Accept		json
// @Produce		json
// @Param			id	path		int	true	"Recipe ID"
// @Success		204
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/recipe/{id}/ingredients [delete]
func (e RecipeEndpoints) DeleteIngredientLink(ctx *gin.Context) {
	e.handlers.DeleteIngredientLink(ctx.Writer, ctx.Request, ctx.Param("id"))
}
