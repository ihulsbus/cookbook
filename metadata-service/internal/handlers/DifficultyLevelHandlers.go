package handlers

import (
	"net/http"

	m "metadata-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DifficultyLevelService interface {
	FindAll() ([]m.DifficultyLevelDTO, error)
	FindSingle(difficultyLevelDTO m.DifficultyLevelDTO) (m.DifficultyLevelDTO, error)
	Create(difficultyLevelDTO m.DifficultyLevelDTO) (m.DifficultyLevelDTO, error)
	Update(difficultyLevelDTO m.DifficultyLevelDTO) (m.DifficultyLevelDTO, error)
	Delete(difficultyLevelDTO m.DifficultyLevelDTO) error
}

type DifficultyLevelHandlers struct {
	difficultyLevelService DifficultyLevelService
	logger                 m.LoggerInterface
}

func NewDifficultyLevelHandlers(difficultyLevels DifficultyLevelService, logger m.LoggerInterface) *DifficultyLevelHandlers {
	return &DifficultyLevelHandlers{
		difficultyLevelService: difficultyLevels,
		logger:                 logger,
	}
}

func (h *DifficultyLevelHandlers) GetAll(ctx *gin.Context) {
	difficultyLevelDTO, err := h.difficultyLevelService.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no difficultyLevels found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, difficultyLevelDTO)
}

func (h *DifficultyLevelHandlers) Get(ctx *gin.Context) {
	var difficultyLevelDTO m.DifficultyLevelDTO
	var err error

	difficultyLevelDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid difficultyLevel ID"})
		return
	}

	difficultyLevelDTO, err = h.difficultyLevelService.FindSingle(difficultyLevelDTO)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "difficultyLevel not found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, difficultyLevelDTO)
}

func (h *DifficultyLevelHandlers) Create(ctx *gin.Context) {
	var difficultyLevelDTO m.DifficultyLevelDTO
	var err error

	if err = ctx.ShouldBindJSON(&difficultyLevelDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	difficultyLevelDTO, err = h.difficultyLevelService.Create(difficultyLevelDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, difficultyLevelDTO)
}

func (h *DifficultyLevelHandlers) Update(ctx *gin.Context) {
	var difficultyLevelDTO m.DifficultyLevelDTO
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid difficultyLevel ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&difficultyLevelDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// deliberaly set this to ensure the parameter ID is used instead of an accidental id in body
	// perhaps separate create/update DTO's are needed
	difficultyLevelDTO.ID = id

	difficultyLevelDTO, err = h.difficultyLevelService.Update(difficultyLevelDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, difficultyLevelDTO)
}

func (h *DifficultyLevelHandlers) Delete(ctx *gin.Context) {
	var difficultyLevelDTO m.DifficultyLevelDTO
	var err error

	difficultyLevelDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid difficultyLevel ID"})
		return
	}

	err = h.difficultyLevelService.Delete(difficultyLevelDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
