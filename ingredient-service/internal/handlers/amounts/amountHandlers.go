package handlers

import (
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
)

type AmountService interface {
	Find(recipeID uuid.UUID) (*[]models.AmountDTO, error)
	Create(recipeID uuid.UUID, amountsDTO *[]models.AmountDTO) (*[]models.AmountDTO, error)
	Update(recipeID uuid.UUID, amountsDTO *[]models.AmountDTO) (*[]models.AmountDTO, error)
	Delete(recipeID uuid.UUID) error
}

type AmountHandlers struct {
	amountService AmountService
	logger        m.LoggerInterface
}

func NewAmountHandlers(ingredients AmountService, logger m.LoggerInterface) *AmountHandlers {
	return &AmountHandlers{
		amountService: ingredients,
		logger:        logger,
	}
}

func (h AmountHandlers) Find(ctx *gin.Context) {
	var recipeID uuid.UUID
	var amountDTO *[]models.AmountDTO
	var err error

	recipeID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	amountDTO, err = h.amountService.Find(recipeID)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no ingredients for recipe found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, amountDTO)
}

func (h AmountHandlers) Create(ctx *gin.Context) {
	var recipeID uuid.UUID
	var amountDTO []models.AmountDTO
	var err error

	recipeID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&amountDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amountDTOResponse, err := h.amountService.Create(recipeID, &amountDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, amountDTOResponse)

}

func (h AmountHandlers) Update(ctx *gin.Context) {
	var recipeID uuid.UUID
	var amountDTO []models.AmountDTO
	var err error

	recipeID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&amountDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amountDTOResponse, err := h.amountService.Update(recipeID, &amountDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, amountDTOResponse)
}

func (h AmountHandlers) Delete(ctx *gin.Context) {
	var recipeID uuid.UUID
	var err error

	recipeID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	err = h.amountService.Delete(recipeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
