package authz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
)

// HTTPDoer is kept local (rather than depending on a concrete client
// type) so any client with a Do method - e.g.
// shared/httpclient.EnhancedHTTPClient - can be used here.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewRecipeServiceFallback builds a MissFallback that asks recipe-service's
// internal access endpoint for the authoritative RecipeAccessState on a
// projection miss, and upserts the result into projector before the
// Authorizer retries its lookup.
//
// client should be a shared/httpclient.EnhancedHTTPClient (retry + circuit
// breaker) built via NewInternalHTTPClient, since this network-internal
// call skips the OAuth client-credentials flow. baseURL is recipe-service's
// internal (non-ingress) base URL.
func NewRecipeServiceFallback(client HTTPDoer, baseURL string, projector *Projector) MissFallback {
	return func(ctx context.Context, recipeID uuid.UUID) error {
		url := fmt.Sprintf("%s/internal/v1/recipe/%s/access", baseURL, recipeID)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("authz: failed to build internal access request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("authz: internal access request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("authz: recipe-service has no record of recipe %s", recipeID)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("authz: unexpected status %d from recipe-service internal endpoint", resp.StatusCode)
		}

		var state rmq.RecipeAccessState
		if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
			return fmt.Errorf("authz: failed to decode internal access response: %w", err)
		}

		return projector.Handle(state)
	}
}
