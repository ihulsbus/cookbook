package handlers

import (
	"net/http"

	m "metadata-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PreparationTimeService interface {
	FindAll() ([]m.PreparationTimeDTO, error)
	FindSingle(preparationTimeDTO m.PreparationTimeDTO) (m.PreparationTimeDTO, error)
	Create(preparationTimeDTO m.PreparationTimeDTO) (m.PreparationTimeDTO, error)
	Update(preparationTimeDTO m.PreparationTimeDTO) (m.PreparationTimeDTO, error)
	Delete(preparationTimeDTO m.PreparationTimeDTO) error
}

type PreparationTimeHandlers struct {
	preparationTimeService PreparationTimeService
	logger                 m.LoggerInterface
}

func NewPreparationTimeHandlers(preparationTimes PreparationTimeService, logger m.LoggerInterface) *PreparationTimeHandlers {
	return &PreparationTimeHandlers{
		preparationTimeService: preparationTimes,
		logger:                 logger,
	}
}

func (h *PreparationTimeHandlers) GetAll(ctx *gin.Context) {
	preparationTimeDTO, err := h.preparationTimeService.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no preparationTimes found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, preparationTimeDTO)
}

func (h *PreparationTimeHandlers) Get(ctx *gin.Context) {
	var preparationTimeDTO m.PreparationTimeDTO
	var err error

	preparationTimeDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid preparationTime ID"})
		return
	}

	preparationTimeDTO, err = h.preparationTimeService.FindSingle(preparationTimeDTO)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "preparationTime not found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, preparationTimeDTO)
}

func (h *PreparationTimeHandlers) Create(ctx *gin.Context) {
	var preparationTimeDTO m.PreparationTimeDTO
	var err error

	if err = ctx.ShouldBindJSON(&preparationTimeDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	preparationTimeDTO, err = h.preparationTimeService.Create(preparationTimeDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, preparationTimeDTO)
}

func (h *PreparationTimeHandlers) Update(ctx *gin.Context) {
	var preparationTimeDTO m.PreparationTimeDTO
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid preparationTime ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&preparationTimeDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// deliberaly set this to ensure the parameter ID is used instead of an accidental id in body
	// perhaps separate create/update DTO's are needed
	preparationTimeDTO.ID = id

	preparationTimeDTO, err = h.preparationTimeService.Update(preparationTimeDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, preparationTimeDTO)
}

func (h *PreparationTimeHandlers) Delete(ctx *gin.Context) {
	var preparationTimeDTO m.PreparationTimeDTO
	var err error

	preparationTimeDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid preparationTime ID"})
		return
	}

	err = h.preparationTimeService.Delete(preparationTimeDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
