package handlers

import (
	m "github.com/ihulsbus/cookbook/shared/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
)

type SearchService interface {
	SearchInstruction(models.InstructionSearchRequestDTO) (models.InstructionSearchResultDTO, error)
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
	var searchRequestDTO models.InstructionSearchRequestDTO
	var err error

	searchRequestDTO.RecipeID = uuid.MustParse(ctx.Query("recipeID"))

	searchResultDTO, err := h.searchService.SearchInstruction(searchRequestDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, searchResultDTO)
}
