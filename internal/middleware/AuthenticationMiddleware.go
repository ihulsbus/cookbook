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
	log "github.com/sirupsen/logrus"
)

var (
	userCtxKey = userContextKey("user")
	claims     struct {
		Name           string `json:"name"`
		Nickname       string `json:"nickname"`
		Picture        string `json:"picture"`
		UpdatedAt      string `json:"updated_at"`
		Email          string `json:"email"`
		Email_verified bool   `json:"email_verified"`
		Iss            string `json:"iss"`
		Sub            string `json:"sub"`
		Aud            string `json:"aud"`
		Exp            int    `json:"exp"`
		Iat            int    `json:"iat"`
		Nonce          string `json:"nonce"`
		At_hash        string `json:"at_hash"`
	}
)

type OidcMW struct {
	oidcConfig *m.OidcConfig
	context    context.Context
	provider   *oidc.Provider
	verifier   *oidc.IDTokenVerifier
	logger     *log.Logger
}

type userContextKey string

func (o *OidcMW) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")
		if token == "" {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}
		bearer := strings.Split(token, " ")
		if len(bearer) != 2 || bearer[0] != "Bearer" {
			ctx.AbortWithError(http.StatusForbidden, errors.New("no valid token found"))
			return
		}

		user, err := o.authorizeUser(bearer[1])
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("bad request"))
			return
		}

		if user != nil {
			o.logger.Debugf("Authenticated user as %s (%s)", user.Username, user.Email)
			// Pass down the request to the next middleware (or final handler)
			ctx.Request = ctx.Request.Clone(o.WithUser(ctx.Request.Context(), user))
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

	return &m.User{Username: claims.Nickname, Name: claims.Name, Avatar: claims.Picture, Email: claims.Email}, nil
}

// WithUser puts the authenticated user information into the current context.
func (o *OidcMW) WithUser(cntx context.Context, authenticatedUser *m.User) context.Context {
	return context.WithValue(cntx, userCtxKey, authenticatedUser)
}

// UserFromContext retrieves information about the authenticated user from the context of the request.
func (o *OidcMW) UserFromContext(ctx context.Context) (*m.User, error) {
	v := ctx.Value(userCtxKey)

	if v == nil {
		return nil, errors.New("no authenticated user found in context")
	}
	return v.(*m.User), nil
}
