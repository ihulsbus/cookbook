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

func TestRecipeQuery_OK(t *testing.T) {
	recipeID := uuid.New()
	var gotAuth string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		assert.Equal(t, "/recipe", r.URL.Path)
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "10", r.URL.Query().Get("limit"))

		page := m.PaginatedResponse[m.RecipeDTO]{
			Data: []m.RecipeDTO{{ID: recipeID, Name: "apple pie", Description: "pie", ServingCount: 4}},
		}
		json.NewEncoder(w).Encode(page)
	}))
	defer srv.Close()

	r := NewRecipeRepository(srv.URL)
	results, err := r.Query(context.Background(), m.SearchRequest{Query: "pie", Page: 2, Limit: 10}, "Bearer test-token")

	assert.NoError(t, err)
	assert.Equal(t, "Bearer test-token", gotAuth)
	assert.Len(t, results, 1)
	assert.Equal(t, recipeID, results[0].ID)
	assert.Equal(t, "apple pie", results[0].Name)
}

func TestRecipeQuery_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	r := NewRecipeRepository(srv.URL)
	_, err := r.Query(context.Background(), m.SearchRequest{Query: "pie"}, "Bearer bad-token")

	assert.Error(t, err)
}

func TestRecipeQuery_InvalidBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	r := NewRecipeRepository(srv.URL)
	_, err := r.Query(context.Background(), m.SearchRequest{Query: "pie"}, "Bearer test-token")

	assert.Error(t, err)
}
