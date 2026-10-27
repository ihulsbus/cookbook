package handlers

import (
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
)

type InstructionService interface {
	Find(recipeID uuid.UUID) (*[]models.InstructionDTO, error)
	Create(entityID uuid.UUID, instructionDTO *[]models.InstructionDTO) (*[]models.InstructionDTO, error)
	Update(entityID uuid.UUID, instructionDTO *[]models.InstructionDTO) (*[]models.InstructionDTO, error)
	Delete(recipeID uuid.UUID) error
}

type InstructionHandlers struct {
	instructionService InstructionService
	logger             m.LoggerInterface
}

func NewInstructionHandlers(service InstructionService, logger m.LoggerInterface) *InstructionHandlers {
	return &InstructionHandlers{
		instructionService: service,
		logger:             logger,
	}
}

func (h InstructionHandlers) Get(ctx *gin.Context) {
	var instructionDTO *[]models.InstructionDTO

	entityID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid instruction ID"})
		return
	}

	instructionDTO, err = h.instructionService.Find(entityID)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no instruction found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, instructionDTO)
}

func (h InstructionHandlers) Create(ctx *gin.Context) {
	var instructionDTO []models.InstructionDTO

	entityID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid entity id provided"})
		return
	}

	if entityID == uuid.Nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "entity id is required"})
		return
	}

	if err = ctx.ShouldBindJSON(&instructionDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	instructionDTOResponse, err := h.instructionService.Create(entityID, &instructionDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, instructionDTOResponse)
}

func (h InstructionHandlers) Update(ctx *gin.Context) {
	var instructionDTO []models.InstructionDTO

	entityID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid entity id provided"})
		return
	}

	if err = ctx.ShouldBindJSON(&instructionDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	instructionDTOResponse, err := h.instructionService.Update(entityID, &instructionDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, instructionDTOResponse)
}

func (h InstructionHandlers) Delete(ctx *gin.Context) {

	entityID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid instruction ID"})
		return
	}

	err = h.instructionService.Delete(entityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
