package repositories

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
)

func TestMetadataQuery_OK(t *testing.T) {
	recipeID := uuid.New()
	var gotAuth, gotPath string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path

		dtos := []m.MetadataSearchResultDTO{{RecipeID: recipeID}}
		json.NewEncoder(w).Encode(dtos)
	}))
	defer srv.Close()

	r := NewMetadataRepository(srv.URL)
	results, err := r.Query(context.Background(), m.SearchRequest{Query: "anything"}, "Bearer test-token")

	assert.NoError(t, err)
	assert.Equal(t, "Bearer test-token", gotAuth)
	// metadata-service has no free-text search (KIT-30) - this always
	// hits the privacy-filtered "all" listing, never a query-shaped path.
	assert.Equal(t, "/metadata/search/all", gotPath)
	assert.Len(t, results, 1)
	assert.Equal(t, recipeID, results[0].RecipeID)
}

func TestMetadataQuery_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	r := NewMetadataRepository(srv.URL)
	_, err := r.Query(context.Background(), m.SearchRequest{Query: "anything"}, "Bearer test-token")

	assert.Error(t, err)
}
