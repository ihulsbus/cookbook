package authz

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"gorm.io/gorm"

	"github.com/ihulsbus/cookbook/shared/oidc"
)

// Permission is the access level a request needs on a recipe.
type Permission string

const (
	PermissionRead  Permission = "read"
	PermissionWrite Permission = "write"

	visibilityPublic = "public"
)

// ErrProjectionMiss means no projection row exists yet for a recipe, e.g.
// the request landed right after recipe.created but before the event was
// consumed. See SetMissFallback.
var ErrProjectionMiss = errors.New("authz: no projection for recipe")

type LoggerInterface interface {
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// MissFallback is invoked on a projection miss so the caller can fetch the
// authoritative access state and feed it to a Projector before the lookup
// is retried once. An error leaves the miss unresolved.
type MissFallback func(ctx context.Context, recipeID uuid.UUID) error

// Authorizer answers CanRead/CanWrite against the local RecipeAccess
// projection.
type Authorizer struct {
	db     *gorm.DB
	logger LoggerInterface
	onMiss MissFallback

	missCount uint64
	hitCount  uint64
}

// NewAuthorizer creates an Authorizer backed by db. Without a
// MissFallback (see SetMissFallback), a projection miss is simply denied.
func NewAuthorizer(db *gorm.DB, logger LoggerInterface) *Authorizer {
	return &Authorizer{db: db, logger: logger}
}

// SetMissFallback installs f, called once per projection miss before
// failing the request. Not set by this package itself.
func (a *Authorizer) SetMissFallback(f MissFallback) {
	a.onMiss = f
}

// MissCount returns how many projection misses have been observed. Wire
// this into metrics if a sustained non-zero rate needs alerting.
func (a *Authorizer) MissCount() uint64 { return atomic.LoadUint64(&a.missCount) }

func (a *Authorizer) HitCount() uint64 { return atomic.LoadUint64(&a.hitCount) }

func (a *Authorizer) CanRead(ctx context.Context, principal *oidc.Principal, recipeID uuid.UUID) (bool, error) {
	row, err := a.lookup(ctx, recipeID)
	if err != nil {
		return false, err
	}
	allowed, err := evaluate(row, principal, PermissionRead)
	return allowed, err
}

func (a *Authorizer) CanWrite(ctx context.Context, principal *oidc.Principal, recipeID uuid.UUID) (bool, error) {
	row, err := a.lookup(ctx, recipeID)
	if err != nil {
		return false, err
	}
	allowed, err := evaluate(row, principal, PermissionWrite)
	return allowed, err
}

// lookup fetches the projection row for recipeID, invoking the miss
// fallback (if any) and retrying once on ErrProjectionMiss.
func (a *Authorizer) lookup(ctx context.Context, recipeID uuid.UUID) (*RecipeAccess, error) {
	row, err := a.fetch(ctx, recipeID)
	if err == nil {
		atomic.AddUint64(&a.hitCount, 1)
		return row, nil
	}
	if !errors.Is(err, ErrProjectionMiss) {
		return nil, err
	}

	atomic.AddUint64(&a.missCount, 1)
	misses := atomic.LoadUint64(&a.missCount)
	a.logger.Warnf("authz: projection miss for recipe %s (total misses so far: %d)", recipeID, misses)

	if a.onMiss == nil {
		return nil, err
	}

	if fbErr := a.onMiss(ctx, recipeID); fbErr != nil {
		a.logger.Errorf("authz: projection-miss fallback failed for recipe %s: %v", recipeID, fbErr)
		return nil, err
	}

	return a.fetch(ctx, recipeID)
}

func (a *Authorizer) fetch(ctx context.Context, recipeID uuid.UUID) (*RecipeAccess, error) {
	var row RecipeAccess

	err := a.db.WithContext(ctx).First(&row, "recipe_id = ?", recipeID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectionMiss
		}
		return nil, err
	}
	if row.Deleted {
		return nil, errors.New("authz: recipe deleted")
	}

	return &row, nil
}

// evaluate is the rule matrix as a pure function (no DB, no fallback) for
// table-testability: owner gets read/write; a share targeting the
// principal (by user or group) grants its permission; public visibility
// grants read to anyone. Everything else is denied.
func evaluate(row *RecipeAccess, principal *oidc.Principal, perm Permission) (bool, error) {
	if row.OwnerID == principal.UserID {
		return true, nil
	}

	if perm == PermissionRead && row.Visibility == visibilityPublic {
		return true, nil
	}

	var shares []rmq.ShareEntry
	if len(row.Shares) > 0 {
		if err := json.Unmarshal(row.Shares, &shares); err != nil {
			return false, err
		}
	}

	for _, share := range shares {
		var targetsPrincipal bool
		switch share.TargetType {
		case "user":
			targetsPrincipal = share.TargetID == principal.UserID
		case "group":
			targetsPrincipal = containsString(principal.Groups, share.TargetID)
		}

		if !targetsPrincipal {
			continue
		}
		if perm == PermissionRead || share.Permission == string(PermissionWrite) {
			return true, nil
		}
	}

	return false, nil
}

func containsString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// RequireRecipeAccess is a Gin middleware requiring perm on the recipe
// identified by path parameter param. Must run after
// oidc.Module.RequireAuth, which stores the Principal it reads.
func (a *Authorizer) RequireRecipeAccess(param string, perm Permission) gin.HandlerFunc {
	return a.requireRecipeAccess(perm, func(c *gin.Context) (string, bool) {
		v := c.Param(param)
		return v, v != ""
	})
}

// RequireRecipeAccessQuery is RequireRecipeAccess for a recipe id carried
// as a query parameter instead of a path parameter.
func (a *Authorizer) RequireRecipeAccessQuery(param string, perm Permission) gin.HandlerFunc {
	return a.requireRecipeAccess(perm, func(c *gin.Context) (string, bool) {
		return c.GetQuery(param)
	})
}

func (a *Authorizer) requireRecipeAccess(perm Permission, extractID func(*gin.Context) (string, bool)) gin.HandlerFunc {
	return func(c *gin.Context) {
		principal, ok := oidc.PrincipalFromContext(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated principal"})
			return
		}

		raw, ok := extractID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "recipe id is required"})
			return
		}

		recipeID, err := uuid.Parse(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid recipe id"})
			return
		}

		var allowed bool
		switch perm {
		case PermissionRead:
			allowed, err = a.CanRead(c.Request.Context(), principal, recipeID)
		case PermissionWrite:
			allowed, err = a.CanWrite(c.Request.Context(), principal, recipeID)
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "unknown permission"})
			return
		}

		if err != nil {
			if errors.Is(err, ErrProjectionMiss) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "recipe access is not yet known - try again shortly"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "authorization check failed"})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}
