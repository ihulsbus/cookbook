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
	FindAll(pagination models.PaginationRequest, requesterID uuid.UUID, requesterGroups []string, scope models.RecipeListScope) (models.PaginatedResponse[models.RecipeDTO], error)
	FindSingle(recipe models.RecipeDTO, requesterID uuid.UUID, requesterGroups []string) (models.RecipeDTO, error)
	Create(recipe models.RecipeDTO, ownerID uuid.UUID) (models.RecipeDTO, error)
	Update(recipe models.RecipeDTO, requesterID uuid.UUID, requesterGroups []string) (models.RecipeDTO, error)
	Delete(recipe models.RecipeDTO, requesterID uuid.UUID) error
	ListShares(recipeID uuid.UUID, requesterID uuid.UUID) ([]models.RecipeShareDTO, error)
	UpsertShare(recipeID uuid.UUID, requesterID uuid.UUID, share models.RecipeShareDTO) (models.RecipeShareDTO, error)
	DeleteShare(recipeID uuid.UUID, requesterID uuid.UUID, targetType models.ShareTargetType, targetID string) error
	UpdateVisibility(recipeID uuid.UUID, requesterID uuid.UUID, visibility models.Visibility) (models.RecipeDTO, error)
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

// principal extracts the authenticated caller's user ID and group
// memberships from the Gin context (set by oidc.Module.RequireAuth, which
// every route in this service's router already runs), parsing the user ID
// as a UUID. It only fails when RequireAuth wasn't actually run upstream,
// or the OIDC provider's `sub` claim isn't a UUID - both treated as 401,
// since there's no valid caller identity to authorize against.
func (h HttpHandlers) principal(ctx *gin.Context) (uuid.UUID, []string, bool) {
	p, ok := oidc.PrincipalFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated principal"})
		return uuid.Nil, nil, false
	}

	id, err := uuid.Parse(p.UserID)
	if err != nil {
		h.logger.Errorf("principal sub %q is not a UUID: %s", p.UserID, err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid principal"})
		return uuid.Nil, nil, false
	}

	return id, p.Groups, true
}

// principalID is principal without the group memberships, for handlers
// that don't need them (anything that isn't a read/write access check).
func (h HttpHandlers) principalID(ctx *gin.Context) (uuid.UUID, bool) {
	id, _, ok := h.principal(ctx)
	return id, ok
}

// GetAll serves GET /recipe. ?scope=mine|shared|all (KIT-30/KIT-35) narrows
// the result; omitted defaults to "all" (owned + public), the behavior
// this endpoint always had before ?scope existed.
func (h HttpHandlers) GetAll(ctx *gin.Context) {
	requesterID, requesterGroups, ok := h.principal(ctx)
	if !ok {
		return
	}

	pagination, err := hh.ParsePagination(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scope := models.RecipeListScope(ctx.Query("scope"))

	recipeDTO, err := h.recipeService.FindAll(pagination, requesterID, requesterGroups, scope)
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
	requesterID, requesterGroups, ok := h.principal(ctx)
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

	recipeDTO, err = h.recipeService.FindSingle(recipeDTO, requesterID, requesterGroups)
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
	requesterID, requesterGroups, ok := h.principal(ctx)
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

	recipeDTO, err = h.recipeService.Update(recipeDTO, requesterID, requesterGroups)
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

// GetShares serves GET /recipe/:id/shares (KIT-32) - the recipe's owner
// only; see RecipeService.ListShares for why a non-owner gets the same
// "not found" FindSingle uses rather than a distinct permission error.
func (h HttpHandlers) GetShares(ctx *gin.Context) {
	requesterID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	recipeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	shares, err := h.recipeService.ListShares(recipeID, requesterID)
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

	ctx.JSON(http.StatusOK, shares)
}

// UpsertShare serves PUT /recipe/:id/shares (KIT-32): creates a share, or
// updates its permission if one already exists for the same target.
func (h HttpHandlers) UpsertShare(ctx *gin.Context) {
	requesterID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	recipeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	var shareDTO models.RecipeShareDTO
	if err := ctx.ShouldBindJSON(&shareDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	shareDTO, err = h.recipeService.UpsertShare(recipeID, requesterID, shareDTO)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		case "invalid target type", "invalid permission", "target id is required", "target id must be a valid user id":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, shareDTO)
}

// DeleteShare serves DELETE /recipe/:id/shares/:targetType/:targetId
// (KIT-32).
func (h HttpHandlers) DeleteShare(ctx *gin.Context) {
	requesterID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	recipeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	targetType := models.ShareTargetType(ctx.Param("targetType"))
	if !targetType.Valid() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid target type"})
		return
	}

	targetID := ctx.Param("targetId")
	if targetID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "target id is required"})
		return
	}

	err = h.recipeService.DeleteShare(recipeID, requesterID, targetType, targetID)
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

// visibilityRequest is PUT /recipe/:id/visibility's body (KIT-32).
type visibilityRequest struct {
	Visibility models.Visibility `json:"visibility"`
}

// UpdateVisibility serves PUT /recipe/:id/visibility (KIT-32) - a
// lighter-weight alternative to a full recipe Update when only the
// visibility needs to change.
func (h HttpHandlers) UpdateVisibility(ctx *gin.Context) {
	requesterID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	recipeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	var body visibilityRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	recipeDTO, err := h.recipeService.UpdateVisibility(recipeID, requesterID, body.Visibility)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		case "invalid visibility":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, recipeDTO)
}
