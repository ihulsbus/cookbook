package handlers

import (
	"context"
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/gin-gonic/gin"
)

type SearchService interface {
	Search(ctx context.Context, req m.SearchRequest, token string) (m.SearchResult, error)
}

type SearchHandlers struct {
	service SearchService
}

func NewSearchHandlers(service SearchService) *SearchHandlers {
	return &SearchHandlers{
		service: service,
	}
}

func (h *SearchHandlers) Search(ctx *gin.Context) {
	var searchRequest m.SearchRequest
	var searchResult m.SearchResult
	var err error

	if err = ctx.ShouldBindJSON(&searchRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Forwarded as-is to each downstream service (KIT-30), so their own
	// authorization/visibility filtering applies - this service does no
	// filtering of its own. RequireAuth already ran, so this header is
	// present and valid.
	token := ctx.GetHeader("Authorization")

	searchResult, err = h.service.Search(ctx, searchRequest, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, searchResult)
}
