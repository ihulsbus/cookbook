package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandlers interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, categoryID string)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request, categoryID string)
	Delete(w http.ResponseWriter, r *http.Request, categoryID string)
}
type CategoryEndpoints struct {
	handlers CategoryHandlers
}

func NewCategoryEndpoints(handers CategoryHandlers) *CategoryEndpoints {
	return &CategoryEndpoints{
		handlers: handers,
	}
}

// @Summary		Get a list of all available categorys
// @Description	Returns a JSON array of all available categorys
// @Categorys			categorys
// @Produce		json
// @Success		200	{array}		models.Category
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/category [get]
func (e CategoryEndpoints) GetAll(ctx *gin.Context) {
	e.handlers.GetAll(ctx.Writer, ctx.Request)
}

// @Summary		Get a single category
// @Description	Returns the JSON object of a single category
// @Categorys			categorys
// @Produce		json
// @Param			id	path		int	true	"category ID"
// @Success		200	{object}	models.Category
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/category/{id} [get]
func (e CategoryEndpoints) GetSingle(ctx *gin.Context) {
	e.handlers.Get(ctx.Writer, ctx.Request, ctx.Param("id"))
}

// @Summary		Create a new category
// @Description	Creates a new category and returns the JSON object of the created category
// @Categorys			categorys
// @Accept		json
// @Produce		json
// @Param			requestbody	body		models.Category	true	"Create a category"
// @Success		200	{object}	models.Category
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/category [post]
func (e CategoryEndpoints) Create(ctx *gin.Context) {
	e.handlers.Create(ctx.Writer, ctx.Request)
}

// @Summary		Updates an existing category
// @Description	Updates an existing category and returns the JSON object of the updated category
// @Categorys			categorys
// @Accept		json
// @Produce		json
// @Param			id	path		int	true	"category ID"
// @Param			requestbody	body		models.Category	true	"Update a category"
// @Success		200	{object}	models.Category
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/category/{id} [put]
func (e CategoryEndpoints) Update(ctx *gin.Context) {
	e.handlers.Update(ctx.Writer, ctx.Request, ctx.Param("id"))
}

// @Summary		Deletes a category
// @Description	Delete an existing category and returns a simple HTTP code
// @Categorys			categorys
// @Produce		json
// @Param			id	path	int	true	"category ID"
// @Success		204
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/category/{id} [delete]
func (e CategoryEndpoints) Delete(ctx *gin.Context) {
	e.handlers.Delete(ctx.Writer, ctx.Request, ctx.Param("id"))
}
