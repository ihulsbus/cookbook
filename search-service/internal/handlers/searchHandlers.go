package handlers

import (
	"context"
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/gin-gonic/gin"
)

type SearchService interface {
	Search(ctx context.Context, req m.SearchRequest) (m.SearchResult, error)
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

	searchResult, err = h.service.Search(ctx, searchRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, searchResult)
}
