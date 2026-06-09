package handlers

import (
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	hh "github.com/ihulsbus/cookbook/shared/http"
)

const (
	invalidID = "invalid preparationTime ID"
)

type PreparationTimeService interface {
	FindAll(pagination m.PaginationRequest) (m.PaginatedResponse[m.PreparationTime], error)
	FindSingle(preparationTime m.PreparationTime) (m.PreparationTime, error)
	Create(preparationTime m.PreparationTime) (m.PreparationTime, error)
	Update(preparationTime m.PreparationTime) (m.PreparationTime, error)
	Delete(preparationTime m.PreparationTime) error
}

type PreparationTime struct {
	preparationTimeService PreparationTimeService
	logger                 m.LoggerInterface
}

func NewPreparationTimeHandlers(preparationTimes PreparationTimeService, logger m.LoggerInterface) *PreparationTime {
	return &PreparationTime{
		preparationTimeService: preparationTimes,
		logger:                 logger,
	}
}

func (h *PreparationTime) GetAll(ctx *gin.Context) {
	pagination, err := hh.ParsePagination(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	preparationTime, err := h.preparationTimeService.FindAll(pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, preparationTime)
}

func (h *PreparationTime) Get(ctx *gin.Context) {
	var preparationTime m.PreparationTime
	var err error

	preparationTime.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": invalidID})
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

func (h *PreparationTime) Create(ctx *gin.Context) {
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

func (h *PreparationTime) Update(ctx *gin.Context) {
	var preparationTime m.PreparationTime
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": invalidID})
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

func (h *PreparationTime) Delete(ctx *gin.Context) {
	var preparationTime m.PreparationTime
	var err error

	preparationTime.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": invalidID})
		return
	}

	err = h.preparationTimeService.Delete(preparationTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Data(http.StatusNoContent, "application/json", nil)
}
