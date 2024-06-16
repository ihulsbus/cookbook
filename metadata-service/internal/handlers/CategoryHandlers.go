package handlers

import (
	"net/http"

	m "metadata-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryService interface {
	FindAll() ([]m.CategoryDTO, error)
	FindSingle(CategoryDTO m.CategoryDTO) (m.CategoryDTO, error)
	Create(CategoryDTO m.CategoryDTO) (m.CategoryDTO, error)
	Update(CategoryDTO m.CategoryDTO) (m.CategoryDTO, error)
	Delete(CategoryDTO m.CategoryDTO) error
}

type CategoryHandlers struct {
	categoryService CategoryService
	logger          m.LoggerInterface
}

func NewCategoryHandlers(categorys CategoryService, logger m.LoggerInterface) *CategoryHandlers {
	return &CategoryHandlers{
		categoryService: categorys,
		logger:          logger,
	}
}

func (h *CategoryHandlers) GetAll(ctx *gin.Context) {

	categoryDTO, err := h.categoryService.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no categories found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, categoryDTO)
}

func (h *CategoryHandlers) Get(ctx *gin.Context) {
	var categoryDTO m.CategoryDTO
	var err error

	categoryDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	categoryDTO, err = h.categoryService.FindSingle(categoryDTO)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, categoryDTO)
}

func (h *CategoryHandlers) Create(ctx *gin.Context) {
	var categoryDTO m.CategoryDTO
	var err error

	if err = ctx.ShouldBindJSON(&categoryDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categoryDTO, err = h.categoryService.Create(categoryDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, categoryDTO)
}

func (h *CategoryHandlers) Update(ctx *gin.Context) {
	var categoryDTO m.CategoryDTO
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&categoryDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// deliberaly set this to ensure the parameter ID is used instead of an accidental id in body
	// perhaps separate create/update DTO's are needed
	categoryDTO.ID = id

	categoryDTO, err = h.categoryService.Update(categoryDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, categoryDTO)
}

func (h *CategoryHandlers) Delete(ctx *gin.Context) {
	var categoryDTO m.CategoryDTO
	var err error

	categoryDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	err = h.categoryService.Delete(categoryDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
