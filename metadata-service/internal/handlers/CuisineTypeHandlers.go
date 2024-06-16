package handlers

import (
	"net/http"

	m "metadata-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CuisineTypeService interface {
	FindAll() ([]m.CuisineTypeDTO, error)
	FindSingle(cuisineTypeDTO m.CuisineTypeDTO) (m.CuisineTypeDTO, error)
	Create(cuisineTypeDTO m.CuisineTypeDTO) (m.CuisineTypeDTO, error)
	Update(cuisineTypeDTO m.CuisineTypeDTO) (m.CuisineTypeDTO, error)
	Delete(cuisineTypeDTO m.CuisineTypeDTO) error
}

type CuisineTypeHandlers struct {
	cuisineTypeService CuisineTypeService
	logger             m.LoggerInterface
}

func NewCuisineTypeHandlers(cuisineTypes CuisineTypeService, logger m.LoggerInterface) *CuisineTypeHandlers {
	return &CuisineTypeHandlers{
		cuisineTypeService: cuisineTypes,
		logger:             logger,
	}
}

func (h *CuisineTypeHandlers) GetAll(ctx *gin.Context) {
	cuisineTypeDTO, err := h.cuisineTypeService.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no cuisineTypes found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, cuisineTypeDTO)
}

func (h *CuisineTypeHandlers) Get(ctx *gin.Context) {
	var cuisineTypeDTO m.CuisineTypeDTO
	var err error

	cuisineTypeDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid cuisineType ID"})
		return
	}

	cuisineTypeDTO, err = h.cuisineTypeService.FindSingle(cuisineTypeDTO)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "cuisineType not found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, cuisineTypeDTO)
}

func (h *CuisineTypeHandlers) Create(ctx *gin.Context) {
	var cuisineTypeDTO m.CuisineTypeDTO
	var err error

	if err = ctx.ShouldBindJSON(&cuisineTypeDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	cuisineTypeDTO, err = h.cuisineTypeService.Create(cuisineTypeDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, cuisineTypeDTO)
}

func (h *CuisineTypeHandlers) Update(ctx *gin.Context) {
	var cuisineTypeDTO m.CuisineTypeDTO
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid cuisineType ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&cuisineTypeDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// deliberaly set this to ensure the parameter ID is used instead of an accidental id in body
	// perhaps separate create/update DTO's are needed
	cuisineTypeDTO.ID = id

	cuisineTypeDTO, err = h.cuisineTypeService.Update(cuisineTypeDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cuisineTypeDTO)
}

func (h *CuisineTypeHandlers) Delete(ctx *gin.Context) {
	var cuisineTypeDTO m.CuisineTypeDTO
	var err error

	cuisineTypeDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid cuisineType ID"})
		return
	}

	err = h.cuisineTypeService.Delete(cuisineTypeDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
