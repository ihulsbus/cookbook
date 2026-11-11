package oidc

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	testIssuer          = "https://authentik.example.com/application/o/cookbook-web/"
	testAudience        = "cookbook-web"
	testServiceAudience = "cookbook-service"
	testKID             = "test-key"
)

// newTestModule builds a Module backed by a locally generated RSA key and a
// JWKS built from it, bypassing discovery and the network entirely.
func newTestModule(t *testing.T) (*Module, *rsa.PrivateKey) {
	t.Helper()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	jwk, err := jwkset.NewJWKFromKey(&priv.PublicKey, jwkset.JWKOptions{
		Metadata: jwkset.JWKMetadataOptions{KID: testKID, ALG: jwkset.AlgRS256},
	})
	if err != nil {
		t.Fatalf("failed to build JWK: %v", err)
	}

	raw, err := json.Marshal(jwkset.JWKSMarshal{Keys: []jwkset.JWKMarshal{jwk.Marshal()}})
	if err != nil {
		t.Fatalf("failed to marshal JWKS: %v", err)
	}

	keys, err := keyfunc.NewJWKSetJSON(raw)
	if err != nil {
		t.Fatalf("failed to build keyfunc from JWKS: %v", err)
	}

	return newModuleFromKeyfunc(testIssuer, []string{testAudience, testServiceAudience}, keys), priv
}

type signOptions struct {
	issuer   string
	audience string
	groups   []string
	subject  string
	expired  bool
	kid      string
	noKid    bool
}

func signTestToken(t *testing.T, priv *rsa.PrivateKey, opts signOptions) string {
	t.Helper()

	if opts.issuer == "" {
		opts.issuer = testIssuer
	}
	if opts.audience == "" {
		opts.audience = testAudience
	}
	if opts.subject == "" {
		opts.subject = "user-123"
	}
	if opts.kid == "" && !opts.noKid {
		opts.kid = testKID
	}

	expiry := time.Now().Add(time.Hour)
	if opts.expired {
		expiry = time.Now().Add(-time.Hour)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    opts.issuer,
			Subject:   opts.subject,
			Audience:  jwt.ClaimStrings{opts.audience},
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
		PreferredUsername: "jdoe",
		Email:             "jdoe@example.com",
		Groups:            opts.groups,
	})
	if opts.kid != "" {
		token.Header["kid"] = opts.kid
	}

	signed, err := token.SignedString(priv)
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}
	return signed
}

func performRequest(h gin.HandlerFunc, bearer string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", h, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestRequireAuth_MissingHeader(t *testing.T) {
	module, _ := newTestModule(t)

	w := performRequest(module.RequireAuth(), "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuth_InvalidSignature(t *testing.T) {
	module, _ := newTestModule(t)
	_, otherKey := newTestModule(t)

	token := signTestToken(t, otherKey, signOptions{})

	w := performRequest(module.RequireAuth(), token)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuth_UnknownKID(t *testing.T) {
	module, priv := newTestModule(t)

	token := signTestToken(t, priv, signOptions{kid: "does-not-exist"})

	w := performRequest(module.RequireAuth(), token)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuth_Expired(t *testing.T) {
	module, priv := newTestModule(t)

	token := signTestToken(t, priv, signOptions{expired: true})

	w := performRequest(module.RequireAuth(), token)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuth_WrongIssuer(t *testing.T) {
	module, priv := newTestModule(t)

	token := signTestToken(t, priv, signOptions{issuer: "https://not-the-provider.example.com/"})

	w := performRequest(module.RequireAuth(), token)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuth_WrongAudience(t *testing.T) {
	module, priv := newTestModule(t)

	token := signTestToken(t, priv, signOptions{audience: "someone-else"})

	w := performRequest(module.RequireAuth(), token)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuth_Valid(t *testing.T) {
	module, priv := newTestModule(t)

	token := signTestToken(t, priv, signOptions{groups: []string{"cookbook-user"}})

	w := performRequest(module.RequireAuth(), token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// TestRequireAuth_SecondAudience checks that a Module configured with
// multiple accepted audiences accepts tokens for either one.
func TestRequireAuth_SecondAudience(t *testing.T) {
	module, priv := newTestModule(t)

	token := signTestToken(t, priv, signOptions{audience: testServiceAudience})

	w := performRequest(module.RequireAuth(), token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRequireRole_Forbidden(t *testing.T) {
	module, priv := newTestModule(t)

	token := signTestToken(t, priv, signOptions{groups: []string{"cookbook-user"}})

	w := performRequest(module.RequireRole("cookbook-admin"), token)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestRequireRole_Allowed(t *testing.T) {
	module, priv := newTestModule(t)

	token := signTestToken(t, priv, signOptions{groups: []string{"cookbook-user", "cookbook-admin"}})

	w := performRequest(module.RequireRole("cookbook-admin"), token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRequireRole_MissingToken(t *testing.T) {
	module, _ := newTestModule(t)

	w := performRequest(module.RequireRole("cookbook-admin"), "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestPrincipalFromContext(t *testing.T) {
	module, priv := newTestModule(t)
	token := signTestToken(t, priv, signOptions{groups: []string{"cookbook-user"}})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	var gotPrincipal *Principal
	var gotOK bool
	router.GET("/", module.RequireAuth(), func(c *gin.Context) {
		gotPrincipal, gotOK = PrincipalFromContext(c)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if !gotOK {
		t.Fatal("expected principal to be set in context")
	}
	if gotPrincipal.UserID != "user-123" {
		t.Errorf("expected UserID user-123, got %q", gotPrincipal.UserID)
	}
	if gotPrincipal.Username != "jdoe" {
		t.Errorf("expected Username jdoe, got %q", gotPrincipal.Username)
	}
	if len(gotPrincipal.Groups) != 1 || gotPrincipal.Groups[0] != "cookbook-user" {
		t.Errorf("expected Groups [cookbook-user], got %v", gotPrincipal.Groups)
	}
}

// TestDiscover exercises the real discovery path against a local test
// server, independent of the Module tests above which bypass it entirely.
func TestDiscover(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/.well-known/openid-configuration" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(DiscoveryDocument{
			Issuer:  "https://issuer.example.com/",
			JWKSURI: "https://issuer.example.com/jwks.json",
		})
	}))
	defer server.Close()

	doc, err := Discover(t.Context(), server.URL)
	if err != nil {
		t.Fatalf("discover failed: %v", err)
	}
	if doc.JWKSURI != "https://issuer.example.com/jwks.json" {
		t.Errorf("unexpected jwks_uri: %q", doc.JWKSURI)
	}
}

func TestDiscover_MissingJWKSURI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(DiscoveryDocument{Issuer: "https://issuer.example.com/"})
	}))
	defer server.Close()

	_, err := Discover(t.Context(), server.URL)
	if err == nil {
		t.Fatal("expected an error for a discovery document missing jwks_uri")
	}
}
