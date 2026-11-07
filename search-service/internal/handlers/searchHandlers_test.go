package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
)

type searchServiceMock struct {
	gotToken string
}

func (s *searchServiceMock) Search(ctx context.Context, req m.SearchRequest, token string) (m.SearchResult, error) {
	s.gotToken = token
	switch req.Query {
	case "fail":
		return m.SearchResult{}, errors.New("downstream error")
	default:
		return m.SearchResult{Recipes: []m.RecipeResult{{Name: req.Query}}}, nil
	}
}

func TestSearch_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &searchServiceMock{}
	h := NewSearchHandlers(svc)

	reqBody, _ := json.Marshal(m.SearchRequest{Query: "pie"})
	req := httptest.NewRequest("POST", "http://example.com/api/v2/search", bytes.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Search(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Bearer test-token", svc.gotToken)

	var result m.SearchResult
	assert.NoError(t, json.Unmarshal(body, &result))
	assert.Equal(t, "pie", result.Recipes[0].Name)
}

func TestSearch_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewSearchHandlers(&searchServiceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/search", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Search(c)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestSearch_ServiceErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewSearchHandlers(&searchServiceMock{})

	reqBody, _ := json.Marshal(m.SearchRequest{Query: "fail"})
	req := httptest.NewRequest("POST", "http://example.com/api/v2/search", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Search(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"downstream error"}`, string(body))
}
