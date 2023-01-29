package middleware

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	m "github.com/ihulsbus/cookbook/internal/models"
	log "github.com/sirupsen/logrus"
)

type Middleware struct {
	logger *log.Logger
	OidcMW OidcMW
}

// Init a new middleware instance
func NewMiddleware(oidcConfig *m.OidcConfig, logger *log.Logger) *Middleware {
	ctx := context.Background()
	provider, _ := oidc.NewProvider(ctx, oidcConfig.URL)
	verifier := provider.Verifier(&oidc.Config{
		ClientID:             oidcConfig.ClientID,
		SupportedSigningAlgs: oidcConfig.SigningAlgs,
		SkipClientIDCheck:    oidcConfig.SkipClientIDCheck,
		SkipExpiryCheck:      oidcConfig.SkipExpiryCheck,
		SkipIssuerCheck:      oidcConfig.SkipIssuerCheck,
	})

	return &Middleware{
		logger: logger,
		OidcMW: OidcMW{
			context:    ctx,
			provider:   provider,
			verifier:   verifier,
			oidcConfig: oidcConfig,
			logger:     logger,
		},
	}
}
