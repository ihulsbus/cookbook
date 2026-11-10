package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/ihulsbus/cookbook/shared/oidc"

	m "user-service/internal/models"
)

// UserService is the subset of services.UserService's methods this
// handler needs, declared here so it can be unit-tested without the
// real service (and its Authentik client) behind it.
type UserService interface {
	SearchUser(ctx context.Context, query string) (m.UserDTO, error)
	GetDisplayName(ctx context.Context, id uuid.UUID) (m.UserDTO, error)
	ListMyHouseholds(ctx context.Context, userID uuid.UUID) ([]m.HouseholdDTO, error)
	CreateHousehold(ctx context.Context, ownerID uuid.UUID, name string) (m.HouseholdDTO, error)
	AddMember(ctx context.Context, householdID string, requesterID uuid.UUID, query string) (m.HouseholdDTO, error)
	RemoveMember(ctx context.Context, householdID string, requesterID, targetUserID uuid.UUID) error
	LeaveHousehold(ctx context.Context, householdID string, userID uuid.UUID) error
	DeleteHousehold(ctx context.Context, householdID string, requesterID uuid.UUID) error
}

type HttpHandlers struct {
	userService UserService
	logger      *log.Logger
}

func NewHttpHandlers(userService UserService, logger *log.Logger) *HttpHandlers {
	return &HttpHandlers{
		userService: userService,
		logger:      logger,
	}
}

// principalID extracts the authenticated caller's user ID from the Gin
// context (set by oidc.Module.RequireAuth, which every route in this
// service's router already runs). KIT-33/34 never need group
// memberships here - every ownership/membership check is against
// Authentik's own group data, not the token's `groups` claim.
func (h HttpHandlers) principalID(ctx *gin.Context) (uuid.UUID, bool) {
	p, ok := oidc.PrincipalFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated principal"})
		return uuid.Nil, false
	}

	id, err := uuid.Parse(p.UserID)
	if err != nil {
		h.logger.Errorf("principal sub %q is not a UUID: %s", p.UserID, err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid principal"})
		return uuid.Nil, false
	}

	return id, true
}

// SearchUsers serves GET /api/v2/user/users?query= (KIT-33).
func (h HttpHandlers) SearchUsers(ctx *gin.Context) {
	if _, ok := h.principalID(ctx); !ok {
		return
	}

	user, err := h.userService.SearchUser(ctx.Request.Context(), ctx.Query("query"))
	if err != nil {
		switch err.Error() {
		case "query is required":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "user not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// GetUser serves GET /api/v2/user/users/:id (KIT-33).
func (h HttpHandlers) GetUser(ctx *gin.Context) {
	if _, ok := h.principalID(ctx); !ok {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.userService.GetDisplayName(ctx.Request.Context(), id)
	if err != nil {
		switch err.Error() {
		case "user not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// ListMyHouseholds serves GET /api/v2/user/me/households (KIT-33).
func (h HttpHandlers) ListMyHouseholds(ctx *gin.Context) {
	userID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	households, err := h.userService.ListMyHouseholds(ctx.Request.Context(), userID)
	if err != nil {
		switch err.Error() {
		case "user not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, households)
}

type createHouseholdRequest struct {
	Name string `json:"name"`
}

// CreateHousehold serves POST /api/v2/user/households (KIT-34). The
// caller becomes the household's owner.
func (h HttpHandlers) CreateHousehold(ctx *gin.Context) {
	ownerID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	var body createHouseholdRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	household, err := h.userService.CreateHousehold(ctx.Request.Context(), ownerID, body.Name)
	if err != nil {
		switch err.Error() {
		case "name is required":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "user not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, household)
}

type addMemberRequest struct {
	Query string `json:"query"`
}

// AddMember serves POST /api/v2/user/households/:id/members (KIT-34).
// Owner only, enforced in the service layer.
func (h HttpHandlers) AddMember(ctx *gin.Context) {
	requesterID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	householdID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid household id"})
		return
	}

	var body addMemberRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	household, err := h.userService.AddMember(ctx.Request.Context(), householdID.String(), requesterID, body.Query)
	if err != nil {
		switch err.Error() {
		case "query is required":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "household not found"})
		case "user not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case "forbidden":
			ctx.JSON(http.StatusForbidden, gin.H{"error": "only the household owner can do this"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, household)
}

// RemoveMember serves DELETE /api/v2/user/households/:id/members/:userId
// (KIT-34). Owner only, enforced in the service layer.
func (h HttpHandlers) RemoveMember(ctx *gin.Context) {
	requesterID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	householdID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid household id"})
		return
	}

	targetUserID, err := uuid.Parse(ctx.Param("userId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	err = h.userService.RemoveMember(ctx.Request.Context(), householdID.String(), requesterID, targetUserID)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "household not found"})
		case "forbidden":
			ctx.JSON(http.StatusForbidden, gin.H{"error": "only the household owner can do this"})
		case "owner cannot be removed":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

// LeaveHousehold serves POST /api/v2/user/households/:id/leave (KIT-34).
func (h HttpHandlers) LeaveHousehold(ctx *gin.Context) {
	userID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	householdID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid household id"})
		return
	}

	err = h.userService.LeaveHousehold(ctx.Request.Context(), householdID.String(), userID)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "household not found"})
		case "owner cannot leave":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

// DeleteHousehold serves DELETE /api/v2/user/households/:id (KIT-34).
// Owner only, enforced in the service layer.
func (h HttpHandlers) DeleteHousehold(ctx *gin.Context) {
	requesterID, ok := h.principalID(ctx)
	if !ok {
		return
	}

	householdID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid household id"})
		return
	}

	err = h.userService.DeleteHousehold(ctx.Request.Context(), householdID.String(), requesterID)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "household not found"})
		case "forbidden":
			ctx.JSON(http.StatusForbidden, gin.H{"error": "only the household owner can do this"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
