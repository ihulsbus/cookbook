package handlers

import (
	"net/http"

	m "metadata-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TagService interface {
	FindAll() ([]m.TagDTO, error)
	FindSingle(tagDTO m.TagDTO) (m.TagDTO, error)
	Create(tagDTO m.TagDTO) (m.TagDTO, error)
	Update(tagDTO m.TagDTO) (m.TagDTO, error)
	Delete(tagDTO m.TagDTO) error
}

type TagHandlers struct {
	tagService TagService
	logger     m.LoggerInterface
}

func NewTagHandlers(tags TagService, logger m.LoggerInterface) *TagHandlers {
	return &TagHandlers{
		tagService: tags,
		logger:     logger,
	}
}

func (h *TagHandlers) GetAll(ctx *gin.Context) {
	tagDTO, err := h.tagService.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no tags found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, tagDTO)
}

func (h *TagHandlers) Get(ctx *gin.Context) {
	var tagDTO m.TagDTO
	var err error

	tagDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID"})
		return
	}

	tagDTO, err = h.tagService.FindSingle(tagDTO)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, tagDTO)
}

func (h *TagHandlers) Create(ctx *gin.Context) {
	var tagDTO m.TagDTO
	var err error

	if err = ctx.ShouldBindJSON(&tagDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	tagDTO, err = h.tagService.Create(tagDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, tagDTO)
}

func (h *TagHandlers) Update(ctx *gin.Context) {
	var tagDTO m.TagDTO
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&tagDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// deliberaly set this to ensure the parameter ID is used instead of an accidental id in body
	// perhaps separate create/update DTO's are needed
	tagDTO.ID = id

	tagDTO, err = h.tagService.Update(tagDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tagDTO)
}

func (h *TagHandlers) Delete(ctx *gin.Context) {
	var tagDTO m.TagDTO
	var err error

	tagDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID"})
		return
	}

	err = h.tagService.Delete(tagDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
