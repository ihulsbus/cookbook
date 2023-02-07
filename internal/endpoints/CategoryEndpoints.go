package endpoints

import "github.com/gin-gonic/gin"

type CategoryEndpoints struct{}

func NewCategoryEndpoints() *CategoryEndpoints {
	return &CategoryEndpoints{}
}

func (e CategoryEndpoints) NotImplemented(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(501, "not implemented")
}

// @Summary		Get a list of all available tags
// @Description	Returns a JSON array of all available tags
// @Tags			categories
// @Produce		json
// @Success		200	{array}		models.Category
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/category [get]
func (e CategoryEndpoints) GetAll(ctx *gin.Context) {
	e.NotImplemented(ctx)
}

// @Summary		Get a single tag
// @Description	Returns the JSON object of a single tag
// @Tags			categories
// @Produce		json
// @Param			id	path		int	true	"category ID"
// @Success		200	{object}	models.Category
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/category/{id} [get]
func (e CategoryEndpoints) GetSingle(ctx *gin.Context) {
	e.NotImplemented(ctx)
}

// @Summary		Create a new tag
// @Description	Creates a new tag and returns the JSON object of the created tag
// @Tags			categories
// @Accept		json
// @Produce		json
// @Param			requestbody	body		models.Category	true	"Create a category"
// @Success		200	{object}	models.Category
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/category [post]
func (e CategoryEndpoints) Create(ctx *gin.Context) {
	e.NotImplemented(ctx)
}

// @Summary		Updates an existing tag
// @Description	Updates an existing tag and returns the JSON object of the updated tag
// @Tags			categories
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
	e.NotImplemented(ctx)
}

// @Summary		Deletes a tag
// @Description	Delete an existing tag and returns a simple HTTP code
// @Tags			categories
// @Produce		json
// @Success		204
// @Param			id	path		int		true	"category ID"
// @Failure		401	{string}	string	"unauthorized"
// @Failure		404	{string}	string	"not found"
// @Failure		500	{string}	string	"Any error"
// @Router			/category/{id} [delete]
func (e CategoryEndpoints) Delete(ctx *gin.Context) {
	e.NotImplemented(ctx)
}
