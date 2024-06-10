package endpoints

import (
	"net/http"

	m "instruction-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InstructionHandlers interface {
	GetInstruction(w http.ResponseWriter, r *http.Request, recipeID uuid.UUID)
	CreateInstruction(w http.ResponseWriter, r *http.Request, recipeID uuid.UUID)
	UpdateInstruction(w http.ResponseWriter, r *http.Request, recipeID uuid.UUID)
	DeleteInstruction(w http.ResponseWriter, r *http.Request, recipeID uuid.UUID)
}

type Middleware interface {
	UserFromContext(ctx *gin.Context) (*m.User, error)
}

type InstructionEndpoints struct {
	handlers   InstructionHandlers
	middleware Middleware
}

func NewRecipeEndpoints(handlers InstructionHandlers, middleware Middleware) *InstructionEndpoints {
	return &InstructionEndpoints{
		handlers:   handlers,
		middleware: middleware,
	}
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
func (e InstructionEndpoints) GetInstruction(ctx *gin.Context) {
	e.handlers.GetInstruction(ctx.Writer, ctx.Request, uuid.UUID{})
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
func (e InstructionEndpoints) CreateInstruction(ctx *gin.Context) {
	e.handlers.CreateInstruction(ctx.Writer, ctx.Request, uuid.UUID{})
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
func (e InstructionEndpoints) UpdateInstruction(ctx *gin.Context) {
	e.handlers.UpdateInstruction(ctx.Writer, ctx.Request, uuid.UUID{})
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
func (e InstructionEndpoints) DeleteInstruction(ctx *gin.Context) {
	e.handlers.DeleteInstruction(ctx.Writer, ctx.Request, uuid.UUID{})
}
