package handlers

import (
	"net/http"

	m "ingredient-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IngredientService interface {
	FindAll() ([]m.IngredientDTO, error)
	FindUnits() ([]m.UnitDTO, error)
	FindSingle(ingredientDTO m.IngredientDTO) (m.IngredientDTO, error)
	Create(ingredientDTO m.IngredientDTO) (m.IngredientDTO, error)
	Update(ingredientDTO m.IngredientDTO) (m.IngredientDTO, error)
	Delete(ingredientDTO m.IngredientDTO) error
}

type IngredientHandlers struct {
	ingredientService IngredientService
	logger            m.LoggerInterface
}

func NewIngredientHandlers(ingredients IngredientService, logger m.LoggerInterface) *IngredientHandlers {
	return &IngredientHandlers{
		ingredientService: ingredients,
		logger:            logger,
	}
}

// Get all ingredients
func (h IngredientHandlers) GetAll(ctx *gin.Context) {
	var ingredientDTO []m.IngredientDTO
	var err error

	ingredientDTO, err = h.ingredientService.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no ingredients found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, ingredientDTO)
}

// Get all units
func (h IngredientHandlers) GetUnits(ctx *gin.Context) {
	var unitDTO []m.UnitDTO
	var err error

	unitDTO, err = h.ingredientService.FindUnits()
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no units found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, unitDTO)
}

// Get a single ingredient
func (h IngredientHandlers) GetSingle(ctx *gin.Context) {
	var ingredientDTO m.IngredientDTO
	var err error

	ingredientDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ingredient ID"})
		return
	}

	ingredientDTO, err = h.ingredientService.FindSingle(ingredientDTO)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no ingredient found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, ingredientDTO)
}

// Create an ingredient
func (h IngredientHandlers) Create(ctx *gin.Context) {
	var ingredientDTO m.IngredientDTO
	var err error

	if err = ctx.ShouldBindJSON(&ingredientDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	ingredientDTO, err = h.ingredientService.Create(ingredientDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, ingredientDTO)
}

func (h IngredientHandlers) Update(ctx *gin.Context) {
	var ingredientDTO m.IngredientDTO
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ingredient ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&ingredientDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// deliberaly set this to ensure the parameter ID is used instead of an accidental id in body
	// perhaps separate create/update DTO's are needed
	ingredientDTO.ID = id

	ingredientDTO, err = h.ingredientService.Update(ingredientDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ingredientDTO)
}

// Delete an ingredient
func (h IngredientHandlers) Delete(ctx *gin.Context) {
	var ingredientDTO m.IngredientDTO
	var err error

	ingredientDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ingredient ID"})
		return
	}

	err = h.ingredientService.Delete(ingredientDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
