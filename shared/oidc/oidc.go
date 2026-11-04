// Package oidc is a provider-agnostic replacement for shared/keycloak. It
// discovers an OpenID Connect provider's configuration, validates access
// tokens locally against the provider's JWKS (cached, auto-refreshing,
// selected by kid - no gocloak, no introspection round-trip per request),
// and exposes Gin middleware that puts a Principal in the request context.
package oidc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Config configures a Module.
type Config struct {
	// Issuer is the OIDC provider's issuer URL, e.g.
	// https://authentik.example.com/application/o/cookbook-web/. Its
	// .well-known/openid-configuration is used to discover the JWKS URI.
	Issuer string
	// Audiences are the `aud` values this service accepts. A token is
	// valid if its `aud` contains ANY of these. Needs more than one
	// entry when a service accepts both end-user tokens (aud = the
	// frontend's public client ID) and service-to-service tokens (aud =
	// the calling service's own confidential client ID).
	Audiences []string
}

// Principal is the authenticated caller, extracted from a validated access
// token's claims.
type Principal struct {
	UserID   string // sub
	Username string // preferred_username
	Email    string
	Groups   []string
	// Roles mirrors Groups today. Kept as a separate field so a future
	// role-mapping layer (e.g. translating multiple groups into a smaller
	// set of application roles) has somewhere to write distinct values
	// without changing RequireRole's or Principal's shape.
	Roles []string
}

type claims struct {
	jwt.RegisteredClaims
	PreferredUsername string   `json:"preferred_username"`
	Email             string   `json:"email"`
	Groups            []string `json:"groups"`
}

// DiscoveryDocument is the subset of an OIDC provider's
// .well-known/openid-configuration document this module and its sibling
// packages (e.g. shared/httpclient, for its own client-credentials calls)
// care about.
type DiscoveryDocument struct {
	Issuer        string `json:"issuer"`
	JWKSURI       string `json:"jwks_uri"`
	TokenEndpoint string `json:"token_endpoint"`
}

// Module validates access tokens for a single OIDC provider against one or
// more accepted audiences.
type Module struct {
	issuer    string
	audiences []string
	keys      keyfunc.Keyfunc
}

// NewModule discovers the provider at cfg.Issuer and sets up a cached,
// auto-refreshing JWKS client for token validation.
func NewModule(ctx context.Context, cfg Config) (*Module, error) {
	if cfg.Issuer == "" {
		return nil, errors.New("oidc: issuer is required")
	}

	// Drop blanks: callers often build this list from several config
	// fields (e.g. a shared frontend audience plus this service's own,
	// possibly-unset outbound client ID) - an empty string must never
	// become an accepted `aud` value.
	var audiences []string
	for _, a := range cfg.Audiences {
		if a != "" {
			audiences = append(audiences, a)
		}
	}
	if len(audiences) == 0 {
		return nil, errors.New("oidc: at least one non-empty audience is required")
	}

	doc, err := Discover(ctx, cfg.Issuer)
	if err != nil {
		return nil, fmt.Errorf("oidc: discovery failed: %w", err)
	}

	keys, err := keyfunc.NewDefaultCtx(ctx, []string{doc.JWKSURI})
	if err != nil {
		return nil, fmt.Errorf("oidc: failed to set up JWKS client: %w", err)
	}

	return newModuleFromKeyfunc(doc.Issuer, audiences, keys), nil
}

// newModuleFromKeyfunc builds a Module from an already-constructed keyfunc,
// bypassing discovery and HTTP entirely. Used by tests to validate against a
// locally signed JWKS.
func newModuleFromKeyfunc(issuer string, audiences []string, keys keyfunc.Keyfunc) *Module {
	return &Module{
		issuer:    issuer,
		audiences: audiences,
		keys:      keys,
	}
}

// Discover fetches and parses the OIDC provider's
// .well-known/openid-configuration document at issuer.
func Discover(ctx context.Context, issuer string) (*DiscoveryDocument, error) {
	wellKnown := strings.TrimRight(issuer, "/") + "/.well-known/openid-configuration"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, wellKnown, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from %s", resp.StatusCode, wellKnown)
	}

	var doc DiscoveryDocument
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode discovery document: %w", err)
	}
	if doc.JWKSURI == "" {
		return nil, errors.New("discovery document is missing jwks_uri")
	}
	if doc.Issuer == "" {
		doc.Issuer = issuer
	}

	return &doc, nil
}

// parse validates a raw bearer token's signature, issuer, audience and
// expiry, and returns the resulting Principal.
func (m *Module) parse(tokenString string) (*Principal, error) {
	var c claims

	token, err := jwt.ParseWithClaims(tokenString, &c, m.keys.Keyfunc, jwt.WithValidMethods([]string{"RS256"}))
	if err != nil {
		return nil, fmt.Errorf("oidc: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("oidc: invalid token")
	}
	if c.Issuer != m.issuer {
		return nil, errors.New("oidc: unexpected issuer")
	}

	audienceOK := false
	for _, a := range c.Audience {
		for _, allowed := range m.audiences {
			if a == allowed {
				audienceOK = true
				break
			}
		}
	}
	if !audienceOK {
		return nil, errors.New("oidc: audience not allowed")
	}

	return &Principal{
		UserID:   c.Subject,
		Username: c.PreferredUsername,
		Email:    c.Email,
		Groups:   c.Groups,
		Roles:    c.Groups,
	}, nil
}

const principalContextKey = "oidc.principal"

// authenticate validates the request's bearer token, aborting the request
// with 401 on failure. On success it stores the Principal in the gin
// context and returns it.
func (m *Module) authenticate(c *gin.Context) *Principal {
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
		return nil
	}

	principal, err := m.parse(strings.TrimPrefix(authHeader, "Bearer "))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return nil
	}

	c.Set(principalContextKey, principal)
	return principal
}

// RequireAuth is a Gin middleware that requires a valid bearer token, with
// no role check. Any authenticated caller passes.
func (m *Module) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if principal := m.authenticate(c); principal != nil {
			c.Next()
		}
	}
}

// RequireRole is a Gin middleware that requires a valid bearer token whose
// Principal.Roles contains role.
func (m *Module) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		principal := m.authenticate(c)
		if principal == nil {
			return
		}

		for _, r := range principal.Roles {
			if r == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}

// PrincipalFromContext returns the Principal stored by RequireAuth/RequireRole.
func PrincipalFromContext(c *gin.Context) (*Principal, bool) {
	v, ok := c.Get(principalContextKey)
	if !ok {
		return nil, false
	}
	p, ok := v.(*Principal)
	return p, ok
}

// SetPrincipal stores p in c exactly as RequireAuth/RequireRole would.
// Exported for handler tests that need an authenticated context without
// going through a real bearer token.
func SetPrincipal(c *gin.Context, p *Principal) {
	c.Set(principalContextKey, p)
}
