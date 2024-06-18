package handlers

import (
	m "metadata-service/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchService interface {
	SearchMetadata(m.MetadataSearchRequestDTO) ([]m.MetadataSearchResultDTO, error)
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

func (h *SearchHandlers) SearchMetadata(ctx *gin.Context) {
	var searchRequestDTO m.MetadataSearchRequestDTO
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
