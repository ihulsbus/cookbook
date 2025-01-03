package keycloak

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
)

type JWKS struct {
	Keys []struct {
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		Alg string `json:"alg"`
		N   string `json:"n"` // Modulus
		E   string `json:"e"` // Exponent
	} `json:"keys"`
}

type KeycloakModule struct {
	Client       *gocloak.GoCloak
	ClientID     string
	ClientSecret string
	Realm        string
	JWKSURL      string
	PublicKey    *rsa.PublicKey
}

type KeyCloakConfig struct {
	ClientID     string
	ClientSecret string
	Url          string
	Realm        string
}

type ResourceAccess struct {
	Roles []string `json:"roles"`
}

type ResourceAccessMap map[string]ResourceAccess

// NewKeycloakModule initializes the Keycloak module
func NewKeycloakModule(config KeyCloakConfig) (*KeycloakModule, error) {
	jwksURL := config.Url + "/realms/" + config.Realm + "/protocol/openid-connect/certs"
	client := gocloak.NewClient(config.Url)

	module := &KeycloakModule{
		Client:       client,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Realm:        config.Realm,
		JWKSURL:      jwksURL,
	}

	// Fetch the public key during initialization
	publicKey, err := fetchRealmPublicKey(jwksURL)
	if err != nil {
		return nil, err
	}
	module.PublicKey = publicKey

	return module, nil
}

// fetchRealmPublicKey fetches the realm's public key dynamically
func fetchRealmPublicKey(jwksURL string) (*rsa.PublicKey, error) {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch JWKS: unexpected status code")
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	if len(jwks.Keys) == 0 {
		return nil, errors.New("no keys found in JWKS")
	}

	key := jwks.Keys[0] // For simplicity, using the first key. Extend for kid selection.

	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, err
	}

	e := 0
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: e,
	}, nil
}

// ValidateToken validates an OAuth2 token using introspection
func (km *KeycloakModule) ValidateToken(ctx context.Context, token string) (*gocloak.IntroSpectTokenResult, error) {
	introspection, err := km.Client.RetrospectToken(ctx, token, km.ClientID, km.ClientSecret, km.Realm)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return nil, err
	}

	if !*introspection.Active {
		return nil, errors.New("invalid or inactive token")
	}

	return introspection, nil
}

// handleError is a helper function to log and respond with an error
func (km *KeycloakModule) handleError(c *gin.Context, statusCode int, message string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
	}
	c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
}

// Middleware provides a Gin middleware for token validation
func (km *KeycloakModule) Middleware(validAudiences []string, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			km.handleError(c, http.StatusUnauthorized, "missing Authorization header", nil)
			return
		}

		token := authHeader[len("Bearer "):] // Strip the "Bearer " prefix
		introspection, err := km.ValidateToken(c.Request.Context(), token)
		if err != nil {
			km.handleError(c, http.StatusUnauthorized, "invalid token", err)
			return
		}

		if introspection.Aud == nil {
			km.handleError(c, http.StatusUnauthorized, "audience not allowed", nil)
			return
		}

		_, claims, err := km.Client.DecodeAccessToken(c.Request.Context(), token, km.Realm)
		claimsMap := *claims
		if err != nil {
			km.handleError(c, http.StatusUnauthorized, "invalid token", nil)
			return
		}

		valid := false
		for _, a := range *introspection.Aud {
			for _, validAudience := range validAudiences {
				if a == validAudience {
					valid = true
					break
				}
			}
		}

		if !valid {
			km.handleError(c, http.StatusUnauthorized, "audience not allowed", nil)
			return
		}

		// Extract the "resource_access" field
		resourceAccessRaw, ok := claimsMap["resource_access"]
		if !ok {
			km.handleError(c, http.StatusForbidden, "resource_access field not found", nil)
			return
		}

		// Marshal the raw data to JSON
		resourceAccessJSON, err := json.Marshal(resourceAccessRaw)
		if err != nil {
			km.handleError(c, http.StatusInternalServerError, "failed to marshal resource_access", err)
		}

		// Decode JSON into the ResourceAccessMap struct
		var resourceAccess ResourceAccessMap
		if err := json.Unmarshal(resourceAccessJSON, &resourceAccess); err != nil {
			km.handleError(c, http.StatusInternalServerError, "failed to unmarshal resource_access", err)
		}

		if requiredRole != "" {
			roleValid := false
			for _, aud := range validAudiences {
				if accessMap, ok := resourceAccess[aud]; ok {
					for _, role := range accessMap.Roles {
						if role == requiredRole {
							roleValid = true
							break
						}
					}
					if roleValid {
						break
					}
				}
			}

			if !roleValid {
				km.handleError(c, http.StatusForbidden, "insufficient permissions", nil)
				return
			}
		}

		c.Next()
	}
}
