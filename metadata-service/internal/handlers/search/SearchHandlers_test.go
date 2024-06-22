package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
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
	minTime int       = 1
	maxTime int       = 2
	id      uuid.UUID = uuid.New()

	searchRequestDTO m.MetadataSearchRequestDTO = m.MetadataSearchRequestDTO{
		RecipeID:          id,
		CategoryID:        id,
		TagID:             id,
		DifficultyLevelID: id,
		CuisineTypeID:     id,
		MinPrepTime:       &minTime,
		MaxPrepTime:       &maxTime,
	}

	searchResultDTO m.MetadataSearchResultDTO = m.MetadataSearchResultDTO{
		RecipeID:          id,
		CategoryIDs:       []uuid.UUID{id},
		TagIDs:            []uuid.UUID{id},
		DifficultyLevelID: id,
		PreparationTimeID: id,
		CuisineTypeID:     id,
	}
)

// ====== SearchService ======

func (s *SearchServiceMock) SearchMetadata(request m.MetadataSearchRequestDTO) ([]m.MetadataSearchResultDTO, error) {
	switch *request.MinPrepTime {
	case 1:
		var response []m.MetadataSearchResultDTO
		response = append(response, searchResultDTO)
		return response, nil
	default:
		return nil, errors.New("error")
	}
}

// ====== Tests ======

func TestSearch_OK(t *testing.T) {
	h := NewSearchHandlers(&SearchServiceMock{}, &m.LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(searchRequestDTO)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.SearchMetadata(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal([]m.MetadataSearchResultDTO{searchResultDTO})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestSearch_UnmarshalErr(t *testing.T) {
	h := NewSearchHandlers(&SearchServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/tag/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.SearchMetadata(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestSearch_SearchErr(t *testing.T) {
	h := NewSearchHandlers(&SearchServiceMock{}, &m.LoggerInterfaceMock{})

	minTime = 2

	reqBody, _ := json.Marshal(searchRequestDTO)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.SearchMetadata(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}
