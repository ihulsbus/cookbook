package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	m "github.com/ihulsbus/cookbook/internal/models"
	"golang.org/x/exp/slices"
)

var (
	userCtxKey = "user"
	claims     m.Claims
)

type OidcMW struct {
	oidcConfig *m.OidcConfig
	context    context.Context
	provider   *oidc.Provider
	verifier   *oidc.IDTokenVerifier
	logger     LoggingInterface
}

func (o *OidcMW) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")
		if token == "" {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}
		bearer := strings.Split(token, " ")
		if len(bearer) != 2 || bearer[0] != "Bearer" {
			_ = ctx.AbortWithError(http.StatusForbidden, errors.New("no valid token found"))
			return
		}

		user, err := o.authorizeUser(bearer[1])
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("bad request"))
			return
		}

		if user != nil {
			o.logger.Debugf("Authenticated user as %s", user.UserID)
			// Set user properly in Gin Context
			ctx.Set(userCtxKey, user)
			// Pass down the request to the next middleware (or final handler)
			ctx.Next()
		}
	}
}

func (o *OidcMW) authorizeUser(bearer string) (*m.User, error) {
	idToken, err := o.verifier.Verify(o.context, bearer)

	if err != nil {
		err = fmt.Errorf("unable to verify token: %v", bearer)
		return nil, err
	}

	if err = idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to get claims: %v", err)
	}

	if !claims.Email_verified {
		return nil, fmt.Errorf("email not verified: %v", claims.Email)
	}

	return &m.User{UserID: claims.FederatedClaims.UserID, Groups: claims.Groups}, nil
}

// UserFromContext retrieves information about the authenticated user from the context of the request.
func (o *OidcMW) UserFromContext(ctx context.Context) (*m.User, error) {
	v := ctx.Value(userCtxKey)

	if v == nil {
		return nil, errors.New("no authenticated user found in context")
	}
	return v.(*m.User), nil
}

func (a *OidcMW) VerifyAuthorization(service string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		user, err := a.UserFromContext(ctx)
		if err != nil {
			a.logger.Debugf("eror getting user from context: %s", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		action := ctx.Request.Method

		allowedRoles := m.AuthorizationModel[service][action]
		a.logger.Debugf("Service: %v; User: %v; action: %v; groups: %v; allowed: %v;", service, user.UserID, action, user.Groups, allowedRoles)
		for i := range user.Groups {
			index := slices.IndexFunc(allowedRoles, func(r string) bool { return r == user.Groups[i] })

			if index != -1 {
				a.logger.Debugf("user %s has role %s. Passed", user.UserID, allowedRoles[index])
				ctx.Next()
				return
			}
		}

		a.logger.Debugf("user %s failed authorization check. Deny.", user.UserID)
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}
