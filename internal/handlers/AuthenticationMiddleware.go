package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/ihulsbus/cookbook/internal/config"
	m "github.com/ihulsbus/cookbook/internal/models"
	log "github.com/sirupsen/logrus"
)

var (
	userCtxKey  = userContextKey("user")
	ctx         = context.Background()
	provider, _ = oidc.NewProvider(ctx, config.Configuration.Oidc.URL)
	verifier    = provider.Verifier(&oidc.Config{
		ClientID:             config.Configuration.Oidc.ClientID,
		SupportedSigningAlgs: config.Configuration.Oidc.SigningAlgs,
		SkipClientIDCheck:    config.Configuration.Oidc.SkipClientIDCheck,
		SkipExpiryCheck:      config.Configuration.Oidc.SkipExpiryCheck,
		SkipIssuerCheck:      config.Configuration.Oidc.SkipIssuerCheck,
	})

	claims struct {
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
}

type userContextKey string

func (oidcmw *OidcMW) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, "unauthorized. No token supplied", http.StatusUnauthorized)
			return
		}
		bearer := strings.Split(token, " ")
		if len(bearer) != 2 || bearer[0] != "Bearer" {
			http.Error(w, "no valid token found", http.StatusForbidden)
			return
		}

		user, err := authorizeUser(bearer[1])
		if err != nil {
			http.Error(w, fmt.Sprintf("forbidden %v", err.Error()), http.StatusBadRequest)
		}

		if user != nil {
			log.Infof("Authenticated user as %s (%s)", user.Username, user.Email)
			// Pass down the request to the next middleware (or final handler)
			r = r.WithContext(WithUser(r.Context(), user))
			next.ServeHTTP(w, r)
		}

	})
}

func authorizeUser(bearer string) (*m.User, error) {
	idToken, err := verifier.Verify(ctx, bearer)

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
func WithUser(cntx context.Context, authenticatedUser *m.User) context.Context {
	return context.WithValue(cntx, userCtxKey, authenticatedUser)
}

// UserFromContext retrieves information about the authenticated user from the context of the request.
func UserFromContext(ctx context.Context) (*m.User, error) {
	v := ctx.Value(userCtxKey)

	if v == nil {
		return nil, errors.New("no authenticated user found in context")
	}
	return v.(*m.User), nil
}
