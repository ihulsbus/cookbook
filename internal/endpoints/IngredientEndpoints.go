package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IngredientHandlers interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetUnits(w http.ResponseWriter, r *http.Request)
	GetSingle(w http.ResponseWriter, r *http.Request, ingredientID string)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type IngredientEndpoints struct {
	handlers IngredientHandlers
}

func NewIngredientEndpoints(handlers IngredientHandlers) *IngredientEndpoints {
	return &IngredientEndpoints{
		handlers: handlers,
	}
}

func (e IngredientEndpoints) NotImplemented(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(501, "not implemented")
}

// @Summary		Get a list of all available ingredients
// @Description	Returns a JSON array of all available ingredients
// @Tags			ingredients
// @Produce		json
// @Success		200	{array}		models.Ingredient
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/ingredient [get]
func (e IngredientEndpoints) GetAll(ctx *gin.Context) {
	e.handlers.GetAll(ctx.Writer, ctx.Request)
}

// @Summary		Get a list of all available ingredient units
// @Description	Returns a JSON array of all available ingredient units
// @Tags			ingredients
// @Produce		json
// @Success		200	{array}		models.Unit
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/ingredient/units [get]
func (e IngredientEndpoints) GetUnits(ctx *gin.Context) {
	e.handlers.GetUnits(ctx.Writer, ctx.Request)
}

// @Summary		Get a single ingredient
// @Description	Returns a JSON object of a single ingredient
// @Tags			ingredients
// @Produce		json
// @Param			id	path		int	true	"Ingredient ID"
// @Success		200	{object}	models.Ingredient
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/ingredient/{id} [get]
func (e IngredientEndpoints) GetSingle(ctx *gin.Context) {
	e.handlers.GetSingle(ctx.Writer, ctx.Request, ctx.Param("id"))
}

// @Summary		Create an ingredient
// @Description	Create a new ingredient and return its JSON object
// @Tags			ingredients
// @Accept			json
// @Produce		json
// @Success		200	{object}	models.Ingredient
// @Param			requestbody	body		models.Ingredient	true	"Create an ingredient"
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/ingredient [post]
func (e IngredientEndpoints) Create(ctx *gin.Context) {
	e.handlers.Create(ctx.Writer, ctx.Request)
}

// @Summary		Update an Ingredient
// @Description	Update an ingredient and return the updated object's JSON object
// @tags			ingredients
// @Accept			json
// @produce		json
// @Param			id	path		int	true	"Ingredient ID"
// @Param			requestbody	body		models.Ingredient	true	"Update an ingredient"
// @Success		200	{array}		models.Ingredient
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/ingredient/{id} [put]
func (e IngredientEndpoints) Update(ctx *gin.Context) {
	e.handlers.Update(ctx.Writer, ctx.Request)
}

// @Summary		Delete
// @Description	Returns a JSON array of all available ingredients
// @tags			ingredients
// @produce		json
// @Param			id	path	int	true	"Ingredient ID"
// @Success		204
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/ingredient/{id} [delete]
func (e IngredientEndpoints) Delete(ctx *gin.Context) {
	e.handlers.Delete(ctx.Writer, ctx.Request)
}
