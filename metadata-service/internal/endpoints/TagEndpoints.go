package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagHandlers interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, tagID string)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request, tagID string)
	Delete(w http.ResponseWriter, r *http.Request, tagID string)
}
type TagEndpoints struct {
	handlers TagHandlers
}

func NewTagEndpoints(handers TagHandlers) *TagEndpoints {
	return &TagEndpoints{
		handlers: handers,
	}
}

// @Summary		Get a list of all available tags
// @Description	Returns a JSON array of all available tags
// @Tags			tags
// @Produce		json
// @Success		200	{array}		models.Tag
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/tag [get]
func (e TagEndpoints) GetAll(ctx *gin.Context) {
	e.handlers.GetAll(ctx.Writer, ctx.Request)
}

// @Summary		Get a single tag
// @Description	Returns the JSON object of a single tag
// @Tags			tags
// @Produce		json
// @Param			id	path		int	true	"tag ID"
// @Success		200	{object}	models.Tag
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/tag/{id} [get]
func (e TagEndpoints) GetSingle(ctx *gin.Context) {
	e.handlers.Get(ctx.Writer, ctx.Request, ctx.Param("id"))
}

// @Summary		Create a new tag
// @Description	Creates a new tag and returns the JSON object of the created tag
// @Tags			tags
// @Accept		json
// @Produce		json
// @Param			requestbody	body		models.Tag	true	"Create a tag"
// @Success		200	{object}	models.Tag
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/tag [post]
func (e TagEndpoints) Create(ctx *gin.Context) {
	e.handlers.Create(ctx.Writer, ctx.Request)
}

// @Summary		Updates an existing tag
// @Description	Updates an existing tag and returns the JSON object of the updated tag
// @Tags			tags
// @Accept		json
// @Produce		json
// @Param			id	path		int	true	"tag ID"
// @Param			requestbody	body		models.Tag	true	"Update a tag"
// @Success		200	{object}	models.Tag
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/tag/{id} [put]
func (e TagEndpoints) Update(ctx *gin.Context) {
	e.handlers.Update(ctx.Writer, ctx.Request, ctx.Param("id"))
}

// @Summary		Deletes a tag
// @Description	Delete an existing tag and returns a simple HTTP code
// @Tags			tags
// @Produce		json
// @Param			id	path	int	true	"tag ID"
// @Success		204
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/tag/{id} [delete]
func (e TagEndpoints) Delete(ctx *gin.Context) {
	e.handlers.Delete(ctx.Writer, ctx.Request, ctx.Param("id"))
}
