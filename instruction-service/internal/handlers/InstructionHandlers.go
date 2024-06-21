package handlers

import (
	m "instruction-service/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InstructionService interface {
	Find(instruction m.InstructionDTO) (m.InstructionDTO, error)
	Create(instruction m.InstructionDTO) (m.InstructionDTO, error)
	Update(instruction m.InstructionDTO) (m.InstructionDTO, error)
	Delete(instruction m.InstructionDTO) error
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
	var instructionDTO m.InstructionDTO
	var err error

	instructionDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid instruction ID"})
		return
	}

	instructionDTO, err = h.instructionService.Find(instructionDTO)
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
	var instructionDTO m.InstructionDTO
	var err error

	if err = ctx.ShouldBindJSON(&instructionDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	instructionDTO, err = h.instructionService.Create(instructionDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, instructionDTO)
}

func (h InstructionHandlers) Update(ctx *gin.Context) {
	var instructionDTO m.InstructionDTO
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid instruction ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&instructionDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	instructionDTO, err = h.instructionService.Update(instructionDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	instructionDTO.ID = id

	ctx.JSON(http.StatusOK, instructionDTO)
}

func (h InstructionHandlers) Delete(ctx *gin.Context) {
	var instructionDTO m.InstructionDTO
	var err error

	instructionDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid instruction ID"})
		return
	}

	err = h.instructionService.Delete(instructionDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
