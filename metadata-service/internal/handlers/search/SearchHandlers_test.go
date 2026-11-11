package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
	"github.com/ihulsbus/cookbook/shared/oidc"
	"github.com/stretchr/testify/assert"
)

type SearchServiceMock struct {
}

var (
	minTime int       = 1
	maxTime int       = 2
	id      uuid.UUID = uuid.New()

	searchRequestDTO models.MetadataSearchRequestDTO = models.MetadataSearchRequestDTO{
		RecipeID:          id,
		CategoryID:        id,
		TagID:             id,
		DifficultyLevelID: id,
		CuisineTypeID:     id,
		MinPrepTime:       &minTime,
		MaxPrepTime:       &maxTime,
	}

	searchResultDTO models.MetadataSearchResultDTO = models.MetadataSearchResultDTO{
		RecipeID:          id,
		CategoryIDs:       []uuid.UUID{id},
		TagIDs:            []uuid.UUID{id},
		DifficultyLevelID: id,
		PreparationTimeID: id,
		CuisineTypeID:     id,
	}
)

func (s *SearchServiceMock) GetAllRecipeMetadata(requesterID uuid.UUID) (*[]models.MetadataSearchResultDTO, error) {
	return nil, nil
}

func (s *SearchServiceMock) SearchMetadata(request models.MetadataSearchRequestDTO, requesterID uuid.UUID) ([]models.MetadataSearchResultDTO, error) {
	switch *request.MinPrepTime {
	case 1:
		var response []models.MetadataSearchResultDTO
		response = append(response, searchResultDTO)
		return response, nil
	default:
		return nil, errors.New("error")
	}
}

// authedContext builds a gin test context with a principal set, as
// oidc.Module.RequireAuth would before either handler in this package runs.
func authedContext(w *httptest.ResponseRecorder, req *http.Request) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	oidc.SetPrincipal(c, &oidc.Principal{UserID: uuid.New().String()})
	return c
}

func TestSearch_MissingPrincipal(t *testing.T) {
	h := NewSearchHandlers(&SearchServiceMock{}, &models.LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(searchRequestDTO)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w) // no principal set
	c.Request = req

	h.SearchMetadata(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, `{"error":"missing authenticated principal"}`, string(body))
}

func TestSearch_OK(t *testing.T) {
	h := NewSearchHandlers(&SearchServiceMock{}, &models.LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(searchRequestDTO)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.SearchMetadata(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal([]models.MetadataSearchResultDTO{searchResultDTO})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestSearch_UnmarshalErr(t *testing.T) {
	h := NewSearchHandlers(&SearchServiceMock{}, &models.LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/tag/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.SearchMetadata(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestSearch_SearchErr(t *testing.T) {
	h := NewSearchHandlers(&SearchServiceMock{}, &models.LoggerInterfaceMock{})

	minTime = 2

	reqBody, _ := json.Marshal(searchRequestDTO)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.SearchMetadata(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}
