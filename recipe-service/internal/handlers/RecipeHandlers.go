package handlers

import (
	"net/http"

	m "recipe-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecipeService interface {
	FindAll() ([]m.RecipeDTO, error)
	FindSingle(recipe m.RecipeDTO) (m.RecipeDTO, error)
	Create(recipe m.RecipeDTO) (m.RecipeDTO, error)
	Update(recipe m.RecipeDTO) (m.RecipeDTO, error)
	Delete(recipe m.RecipeDTO) error
}

type RecipeHandlers struct {
	recipeService RecipeService
	logger        m.LoggerInterface
}

func NewRecipeHandlers(recipes RecipeService, logger m.LoggerInterface) *RecipeHandlers {
	return &RecipeHandlers{
		recipeService: recipes,
		logger:        logger,
	}
}

func (h RecipeHandlers) GetAll(ctx *gin.Context) {

	recipeDTO, err := h.recipeService.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no recipes found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, recipeDTO)
}

func (h RecipeHandlers) Get(ctx *gin.Context) {
	var recipeDTO m.RecipeDTO
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

func (h RecipeHandlers) Create(ctx *gin.Context) {
	var recipeDTO m.RecipeDTO
	var err error

	if err = ctx.ShouldBindJSON(&recipeDTO); err != nil {
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

func (h RecipeHandlers) Update(ctx *gin.Context) {
	var recipeDTO m.RecipeDTO
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

func (h RecipeHandlers) Delete(ctx *gin.Context) {
	var recipeDTO m.RecipeDTO
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
