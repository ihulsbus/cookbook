package handlers

import (
	m "instruction-service/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SearchService interface {
	SearchInstruction(m.InstructionSearchRequestDTO) (m.InstructionSearchResultDTO, error)
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

func (h *SearchHandlers) SearchInstruction(ctx *gin.Context) {
	var searchRequestDTO m.InstructionSearchRequestDTO
	var err error

	searchRequestDTO.RecipeID = uuid.MustParse(ctx.Query("recipeID"))

	searchResultDTO, err := h.searchService.SearchInstruction(searchRequestDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, searchResultDTO)
}
