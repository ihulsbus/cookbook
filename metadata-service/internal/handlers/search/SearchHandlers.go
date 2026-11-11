package handlers

import (
	m "github.com/ihulsbus/cookbook/shared/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
	"github.com/ihulsbus/cookbook/shared/oidc"
)

type SearchService interface {
	GetAllRecipeMetadata(requesterID uuid.UUID) (*[]models.MetadataSearchResultDTO, error)
	SearchMetadata(request models.MetadataSearchRequestDTO, requesterID uuid.UUID) ([]models.MetadataSearchResultDTO, error)
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

// principalID extracts the caller's user ID (set by oidc.Module.RequireAuth)
// so results can be filtered to recipes they may read; without this, results
// would leak the existence of private recipes to any authenticated caller.
func (h *SearchHandlers) principalID(ctx *gin.Context) (uuid.UUID, bool) {
	principal, ok := oidc.PrincipalFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated principal"})
		return uuid.Nil, false
	}

	id, err := uuid.Parse(principal.UserID)
	if err != nil {
		h.logger.Errorf("principal sub %q is not a UUID: %s", principal.UserID, err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid principal"})
		return uuid.Nil, false
	}

	return id, true
}

func (h *SearchHandlers) GetAllMetadata(ctx *gin.Context) {
	requesterID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	MetadataDTO, err := h.searchService.GetAllRecipeMetadata(requesterID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, MetadataDTO)
}

func (h *SearchHandlers) SearchMetadata(ctx *gin.Context) {
	requesterID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	var searchRequestDTO models.MetadataSearchRequestDTO
	var err error

	if err = ctx.ShouldBindJSON(&searchRequestDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	searchResultDTO, err := h.searchService.SearchMetadata(searchRequestDTO, requesterID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, searchResultDTO)
}
