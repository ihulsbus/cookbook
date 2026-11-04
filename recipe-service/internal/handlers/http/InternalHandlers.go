package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
)

// InternalRecipeService is the subset of RecipeService the internal
// access endpoint needs.
type InternalRecipeService interface {
	GetAccessState(recipeID uuid.UUID) (rmq.RecipeAccessState, error)
}

// InternalHandlers serves recipe-service's internal-only API: endpoints
// meant for other services (shared/authz's projection-miss fallback,
// KIT-28), not for clients. Must be mounted on a router that is never
// exposed through the ingress - see httpInternalServer in
// recipe-service.go, and k8s/values.yaml's "internal" port.
type InternalHandlers struct {
	recipeService InternalRecipeService
	logger        m.LoggerInterface
}

func NewInternalHandlers(recipes InternalRecipeService, logger m.LoggerInterface) *InternalHandlers {
	return &InternalHandlers{
		recipeService: recipes,
		logger:        logger,
	}
}

// GetAccess serves GET /internal/v1/recipe/:id/access, the authoritative
// RecipeAccessState for a recipe - called by another service's
// shared/authz when its own projection hasn't caught up yet (e.g. right
// after recipe.created, before the event has been consumed).
func (h InternalHandlers) GetAccess(ctx *gin.Context) {
	recipeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	state, err := h.recipeService.GetAccessState(recipeID)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		default:
			h.logger.Errorf("internal access lookup failed for recipe %s: %s", recipeID, err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	ctx.JSON(http.StatusOK, state)
}
