package handlers

import (
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	hh "github.com/ihulsbus/cookbook/shared/http"
	"github.com/ihulsbus/cookbook/shared/models"
)

type TagService interface {
	FindAll(pagination models.PaginationRequest) (models.PaginatedResponse[models.TagDTO], error)
	FindSingle(tagDTO models.TagDTO) (models.TagDTO, error)
	Create(tagDTO models.TagDTO) (models.TagDTO, error)
	Update(tagDTO models.TagDTO) (models.TagDTO, error)
	Delete(tagDTO models.TagDTO) error
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
	pagination, err := hh.ParsePagination(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagDTO, err := h.tagService.FindAll(pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tagDTO)
}

func (h *TagHandlers) Get(ctx *gin.Context) {
	var tagDTO models.TagDTO
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
	var tagDTO models.TagDTO
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
	var tagDTO models.TagDTO
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
	var tagDTO models.TagDTO
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
