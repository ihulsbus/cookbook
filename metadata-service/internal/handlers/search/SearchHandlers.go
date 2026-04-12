package handlers

import (
	m "github.com/ihulsbus/cookbook/shared/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ihulsbus/cookbook/shared/models"
)

type SearchService interface {
	GetAllRecipeMetadata() (*[]models.MetadataSearchResultDTO, error)
	SearchMetadata(models.MetadataSearchRequestDTO) ([]models.MetadataSearchResultDTO, error)
}

type SearchHandlers struct {
	searchService SearchService
	logger        m.LoggerInterface
}

func NewSearchHandlers(searchs SearchService, logger m.LoggerInterface) *SearchHandlers {
	return &SearchHandlers{
		searchService: searchs,
		logger:        logger,
	}
}

func (h *SearchHandlers) GetAllMetadata(ctx *gin.Context) {

	MetadataDTO, err := h.searchService.GetAllRecipeMetadata()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, MetadataDTO)
}

func (h *SearchHandlers) SearchMetadata(ctx *gin.Context) {
	var searchRequestDTO models.MetadataSearchRequestDTO
	var err error

	if err = ctx.ShouldBindJSON(&searchRequestDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	searchResultDTO, err := h.searchService.SearchMetadata(searchRequestDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, searchResultDTO)
}
