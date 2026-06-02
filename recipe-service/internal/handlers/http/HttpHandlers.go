package handlers

import (
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	hh "github.com/ihulsbus/cookbook/shared/http"
	"github.com/ihulsbus/cookbook/shared/models"
)

type RecipeService interface {
	FindAll(pagination models.PaginationRequest) (models.PaginatedResponse[models.RecipeDTO], error)
	FindSingle(recipe models.RecipeDTO) (models.RecipeDTO, error)
	Create(recipe models.RecipeDTO) (models.RecipeDTO, error)
	Update(recipe models.RecipeDTO) (models.RecipeDTO, error)
	Delete(recipe models.RecipeDTO) error
}

type HttpHandlers struct {
	recipeService RecipeService
	logger        m.LoggerInterface
}

func NewHttpHandlers(recipes RecipeService, logger m.LoggerInterface) *HttpHandlers {
	return &HttpHandlers{
		recipeService: recipes,
		logger:        logger,
	}
}

func (h HttpHandlers) GetAll(ctx *gin.Context) {
	pagination, err := hh.ParsePagination(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipeDTO, err := h.recipeService.FindAll(pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipeDTO)
}

func (h HttpHandlers) Get(ctx *gin.Context) {
	var recipeDTO models.RecipeDTO
	var err error

	recipeDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	recipeDTO, err = h.recipeService.FindSingle(recipeDTO)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, recipeDTO)
}

func (h HttpHandlers) Create(ctx *gin.Context) {
	var recipeDTO models.RecipeDTO
	var err error

	if err = ctx.ShouldBindJSON(&recipeDTO); err != nil {
		h.logger.Errorf("failed to bind JSON to DTO: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	recipeDTO, err = h.recipeService.Create(recipeDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, recipeDTO)
}

func (h HttpHandlers) Update(ctx *gin.Context) {
	var recipeDTO models.RecipeDTO
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&recipeDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// deliberaly set this to ensure the parameter ID is used instead of an accidental id in body
	// perhaps separate create/update DTO's are needed
	recipeDTO.ID = id

	recipeDTO, err = h.recipeService.Update(recipeDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipeDTO)
}

func (h HttpHandlers) Delete(ctx *gin.Context) {
	var recipeDTO models.RecipeDTO
	var err error

	recipeDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	err = h.recipeService.Delete(recipeDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
