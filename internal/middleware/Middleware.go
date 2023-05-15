package middleware

import (
	"context"

	m "github.com/ihulsbus/cookbook/internal/models"
	log "github.com/sirupsen/logrus"
)

type Middleware struct {
	logger   LoggingInterface
	AuthMW   AuthMW
	ClaimsMW CustomClaims
}

type LoggingInterface interface {
	Debugf(format string, args ...interface{})
	WithFields(fields log.Fields) *log.Entry
}

// Init a new middleware instance
func NewMiddleware(auth0 *m.Auth0Config, logger LoggingInterface) *Middleware {
	ctx := context.Background()

	return &Middleware{
		logger: logger,
		AuthMW: AuthMW{
			auth0:   auth0,
			context: ctx,
			logger:  logger,
		},
	}
}
