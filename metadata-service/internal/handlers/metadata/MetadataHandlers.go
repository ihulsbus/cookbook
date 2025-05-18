package handlers

import (
	"net/http"

	m "metadata-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MetadataService interface {
	FindAll() (*[]m.RecipeMetadataDTO, error)
	Find(recipeID uuid.UUID) (*m.RecipeMetadataDTO, error)
	Create(recipeID uuid.UUID, meta *m.RecipeMetadataDTO) (*m.RecipeMetadataDTO, error)
	Update(recipeID uuid.UUID, meta *m.RecipeMetadataDTO) (*m.RecipeMetadataDTO, error)
	Delete(recipeID uuid.UUID) error
}

type MetadataHandlers struct {
	MetadataService MetadataService
	logger          m.LoggerInterface
}

func NewMetadataHandlers(metadata MetadataService, logger m.LoggerInterface) *MetadataHandlers {
	return &MetadataHandlers{
		MetadataService: metadata,
		logger:          logger,
	}
}

func (h *MetadataHandlers) GetAll(ctx *gin.Context) {

	metadata, err := h.MetadataService.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no metadata for recipes found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, metadata)
}

func (h *MetadataHandlers) Get(ctx *gin.Context) {
	var recipeID uuid.UUID
	var err error

	recipeID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid Recipe ID"})
		return
	}

	metadata, err := h.MetadataService.Find(recipeID)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, metadata)
}

func (h *MetadataHandlers) Create(ctx *gin.Context) {
	var recipeID uuid.UUID
	var metadata m.RecipeMetadataDTO
	var err error

	recipeID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid Recipe ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&metadata); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metadataResponse, err := h.MetadataService.Create(recipeID, &metadata)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, metadataResponse)
}

func (h *MetadataHandlers) Update(ctx *gin.Context) {
	var recipeID uuid.UUID
	var metadata m.RecipeMetadataDTO
	var err error

	recipeID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&metadata); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metadataResponse, err := h.MetadataService.Update(recipeID, &metadata)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, metadataResponse)
}

func (h *MetadataHandlers) Delete(ctx *gin.Context) {
	var recipeID uuid.UUID
	var err error

	recipeID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	err = h.MetadataService.Delete(recipeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
