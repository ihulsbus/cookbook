package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	m "github.com/ihulsbus/cookbook/internal/models"
)

type AuthMW struct {
	auth0   *m.Auth0Config
	context context.Context
	logger  LoggingInterface
}

type CustomClaims struct {
	Permissions []string `json:"permissions"`
}

// Validate does nothing, but we need it to
// satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// HasScope checks whether our claims have a specific scope.
func (c CustomClaims) HasPermission(expectedPermission string) bool {
	for i := range c.Permissions {
		if c.Permissions[i] == expectedPermission {
			return true
		}
	}

	return false
}

func (a *AuthMW) EnsureValidScope(expectedScope string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token := ctx.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		claims := token.CustomClaims.(*CustomClaims)

		if claims.HasPermission(expectedScope) {
			ctx.Next()
			return
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Insufficient scope.")
	}
}

func (a *AuthMW) UserFromContext(ctx *gin.Context) (*m.User, error) {
	var user m.User

	claims := ctx.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	subject := claims.RegisteredClaims.Subject

	if subject == "" {
		return nil, errors.New("user not found in request")
	}

	splitSubject := strings.Split(claims.RegisteredClaims.Subject, "|")

	user.Provider = splitSubject[0]
	user.UserID = splitSubject[1]

	return &m.User{}, nil
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func (a *AuthMW) EnsureValidToken() gin.HandlerFunc {
	issuerURL, err := url.Parse("https://" + a.auth0.Domain + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{a.auth0.Audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return adapter.Wrap(middleware.CheckJWT)
}
