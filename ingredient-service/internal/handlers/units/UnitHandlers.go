package handlers

import (
	"net/http"

	m "ingredient-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UnitService interface {
	FindAll() ([]m.UnitDTO, error)
	FindSingle(unitDTO m.UnitDTO) (m.UnitDTO, error)
	Create(unitDTO m.UnitDTO) (m.UnitDTO, error)
	Update(unitDTO m.UnitDTO) (m.UnitDTO, error)
	Delete(unitDTO m.UnitDTO) error
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

// Get all units
func (h UnitHandlers) GetAll(ctx *gin.Context) {
	var unitDTO []m.UnitDTO
	var err error

	unitDTO, err = h.unitService.FindAll()
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

// Get a single unit
func (h UnitHandlers) GetSingle(ctx *gin.Context) {
	var unitDTO m.UnitDTO
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

// Create an unit
func (h UnitHandlers) Create(ctx *gin.Context) {
	var unitDTO m.UnitDTO
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

func (h UnitHandlers) Update(ctx *gin.Context) {
	var unitDTO m.UnitDTO
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

	// deliberaly set this to ensure the parameter ID is used instead of an accidental id in body
	// perhaps separate create/update DTO's are needed
	unitDTO.ID = id

	unitDTO, err = h.unitService.Update(unitDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, unitDTO)
}

// Delete an unit
func (h UnitHandlers) Delete(ctx *gin.Context) {
	var unitDTO m.UnitDTO
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

	ctx.Status(http.StatusOK)
}
