package handlers

import (
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
)

type UnitService interface {
	FindAll() ([]models.UnitDTO, error)
	FindSingle(unitDTO models.UnitDTO) (models.UnitDTO, error)
	Create(unitDTO models.UnitDTO) (models.UnitDTO, error)
	Update(unitDTO models.UnitDTO) (models.UnitDTO, error)
	Delete(unitDTO models.UnitDTO) error
}

type UnitHandlers struct {
	unitService UnitService
	logger      m.LoggerInterface
}

func NewUnitHandlers(units UnitService, logger m.LoggerInterface) *UnitHandlers {
	return &UnitHandlers{
		unitService: units,
		logger:      logger,
	}
}

// GetAll Get all units
func (h UnitHandlers) GetAll(ctx *gin.Context) {
	var unitDTO []models.UnitDTO
	var err error

	unitDTO, err = h.unitService.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusOK, []models.UnitDTO{})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, unitDTO)
}

// GetSingle Get a single unit
func (h UnitHandlers) GetSingle(ctx *gin.Context) {
	var unitDTO models.UnitDTO
	var err error

	unitDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid unit ID"})
		return
	}

	unitDTO, err = h.unitService.FindSingle(unitDTO)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no unit found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, unitDTO)
}

// Create creates a unit
func (h UnitHandlers) Create(ctx *gin.Context) {
	var unitDTO models.UnitDTO
	var err error

	if err = ctx.ShouldBindJSON(&unitDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	unitDTO, err = h.unitService.Create(unitDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, unitDTO)
}

// Update updates a unit
func (h UnitHandlers) Update(ctx *gin.Context) {
	var unitDTO models.UnitDTO
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid unit ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&unitDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// deliberately set this to ensure the parameter ID is used instead of an accidental id in body
	// perhaps separate create/update DTO's are needed
	unitDTO.ID = id

	unitDTO, err = h.unitService.Update(unitDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, unitDTO)
}

// Delete deletes a unit
func (h UnitHandlers) Delete(ctx *gin.Context) {
	var unitDTO models.UnitDTO
	var err error

	unitDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid unit ID"})
		return
	}

	err = h.unitService.Delete(unitDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
