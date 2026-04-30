package handlers

import (
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PreparationTimeService interface {
	FindAll() ([]m.PreparationTime, error)
	FindSingle(preparationTime m.PreparationTime) (m.PreparationTime, error)
	Create(preparationTime m.PreparationTime) (m.PreparationTime, error)
	Update(preparationTime m.PreparationTime) (m.PreparationTime, error)
	Delete(preparationTime m.PreparationTime) error
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
	preparationTime, err := h.preparationTimeService.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusOK, []m.PreparationTime{})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, preparationTime)
}

func (h *PreparationTimeHandlers) Get(ctx *gin.Context) {
	var preparationTime m.PreparationTime
	var err error

	preparationTime.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid preparationTime ID"})
		return
	}

	preparationTime, err = h.preparationTimeService.FindSingle(preparationTime)
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

	ctx.JSON(http.StatusOK, preparationTime)
}

func (h *PreparationTimeHandlers) Create(ctx *gin.Context) {
	var preparationTime m.PreparationTime
	var err error

	if err = ctx.ShouldBindJSON(&preparationTime); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	preparationTime, err = h.preparationTimeService.Create(preparationTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, preparationTime)
}

func (h *PreparationTimeHandlers) Update(ctx *gin.Context) {
	var preparationTime m.PreparationTime
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid preparationTime ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&preparationTime); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// deliberaly set this to ensure the parameter ID is used instead of an accidental id in body
	// perhaps separate create/update DTO's are needed
	preparationTime.ID = id

	preparationTime, err = h.preparationTimeService.Update(preparationTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, preparationTime)
}

func (h *PreparationTimeHandlers) Delete(ctx *gin.Context) {
	var preparationTime m.PreparationTime
	var err error

	preparationTime.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid preparationTime ID"})
		return
	}

	err = h.preparationTimeService.Delete(preparationTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
