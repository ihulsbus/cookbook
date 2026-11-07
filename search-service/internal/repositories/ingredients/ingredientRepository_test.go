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

func TestIngredientQuery_OK(t *testing.T) {
	ingredientID := uuid.New()
	var gotAuth string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		assert.Equal(t, "/ingredient", r.URL.Path)

		page := m.PaginatedResponse[m.IngredientDTO]{
			Data: []m.IngredientDTO{{ID: ingredientID, Name: "asparagus"}},
		}
		json.NewEncoder(w).Encode(page)
	}))
	defer srv.Close()

	r := NewIngredientRepository(srv.URL)
	results, err := r.Query(context.Background(), m.SearchRequest{Query: "asp", Page: 1, Limit: 25}, "Bearer test-token")

	assert.NoError(t, err)
	assert.Equal(t, "Bearer test-token", gotAuth)
	assert.Len(t, results, 1)
	assert.Equal(t, ingredientID, results[0].ID)
	assert.Equal(t, "asparagus", results[0].Name)
}

func TestIngredientQuery_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	r := NewIngredientRepository(srv.URL)
	_, err := r.Query(context.Background(), m.SearchRequest{Query: "asp"}, "Bearer test-token")

	assert.Error(t, err)
}
