package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	m "metadata-service/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type SearchServiceMock struct {
}

var (
	searchRequestDTO m.MetadataSearchRequestDTO = m.MetadataSearchRequestDTO{
		RecipeID:        uuid.New(),
		CategoryID:      uuid.New(),
		TagID:           uuid.New(),
		DifficultyLevel: 1,
		CuisineType:     "",
		MinPrepTime:     1,
		MaxPrepTime:     2,
	}

	searchResultDTO m.MetadataSearchResultDTO
)

// ====== SearchService ======

func (s *SearchServiceMock) SearchMetadata(m.MetadataSearchRequestDTO) (m.MetadataSearchResultDTO, error) {

	return searchResultDTO, nil
}

// ====== Tests ======

func TestSearch_OK(t *testing.T) {
	h := NewSearchHandlers(&SearchServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(searchRequestDTO)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.SearchMetadata(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(tags)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}
