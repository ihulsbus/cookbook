package handlers

import (
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	hh "github.com/ihulsbus/cookbook/shared/http"
	"github.com/ihulsbus/cookbook/shared/models"
	"github.com/ihulsbus/cookbook/shared/oidc"
)

type RecipeService interface {
	FindAll(pagination models.PaginationRequest, requesterID uuid.UUID, scope models.RecipeListScope) (models.PaginatedResponse[models.RecipeDTO], error)
	FindSingle(recipe models.RecipeDTO, requesterID uuid.UUID) (models.RecipeDTO, error)
	Create(recipe models.RecipeDTO, ownerID uuid.UUID) (models.RecipeDTO, error)
	Update(recipe models.RecipeDTO, requesterID uuid.UUID) (models.RecipeDTO, error)
	Delete(recipe models.RecipeDTO, requesterID uuid.UUID) error
}

type HttpHandlers struct {
	recipeService RecipeService
	logger        m.LoggerInterface
}

func NewHttpHandlers(recipes RecipeService, logger m.LoggerInterface) *HttpHandlers {
	return &HttpHandlers{
		recipeService: recipes,
		logger:        logger,
	}
}

// principalID extracts the authenticated caller's user ID from the Gin
// context (set by oidc.Module.RequireAuth, which every route in this
// service's router already runs) and parses it as a UUID. It only fails
// when RequireAuth wasn't actually run upstream, or the OIDC provider's
// `sub` claim isn't a UUID - both treated as 401, since there's no valid
// caller identity to authorize against.
func (h HttpHandlers) principalID(ctx *gin.Context) (uuid.UUID, bool) {
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

// GetAll serves GET /recipe. ?scope=mine|shared|all (KIT-30) narrows the
// result; omitted defaults to "all" (owned + public), the behavior this
// endpoint always had before ?scope existed.
func (h HttpHandlers) GetAll(ctx *gin.Context) {
	requesterID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	pagination, err := hh.ParsePagination(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scope := models.RecipeListScope(ctx.Query("scope"))

	recipeDTO, err := h.recipeService.FindAll(pagination, requesterID, scope)
	if err != nil {
		switch err.Error() {
		case "invalid scope":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, recipeDTO)
}

func (h HttpHandlers) Get(ctx *gin.Context) {
	requesterID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	var recipeDTO models.RecipeDTO
	var err error

	recipeDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	recipeDTO, err = h.recipeService.FindSingle(recipeDTO, requesterID)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, recipeDTO)
}

func (h HttpHandlers) Create(ctx *gin.Context) {
	ownerID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	var recipeDTO models.RecipeDTO
	var err error

	if err = ctx.ShouldBindJSON(&recipeDTO); err != nil {
		h.logger.Errorf("failed to bind JSON to DTO: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	recipeDTO, err = h.recipeService.Create(recipeDTO, ownerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, recipeDTO)
}

func (h HttpHandlers) Update(ctx *gin.Context) {
	requesterID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	var recipeDTO models.RecipeDTO
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&recipeDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// deliberaly set this to ensure the parameter ID is used instead of an accidental id in body
	// perhaps separate create/update DTO's are needed
	recipeDTO.ID = id

	recipeDTO, err = h.recipeService.Update(recipeDTO, requesterID)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, recipeDTO)
}

func (h HttpHandlers) Delete(ctx *gin.Context) {
	requesterID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	var recipeDTO models.RecipeDTO
	var err error

	recipeDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	err = h.recipeService.Delete(recipeDTO, requesterID)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.Status(http.StatusNoContent)
}
