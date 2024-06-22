package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	m "instruction-service/internal/models"
	"io"
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
	id uuid.UUID = uuid.New()

	searchRequestDTO m.InstructionSearchRequestDTO = m.InstructionSearchRequestDTO{
		RecipeID: id,
	}

	searchResultDTO m.InstructionSearchResultDTO = m.InstructionSearchResultDTO{
		RecipeID:       id,
		InstructionIDs: []uuid.UUID{id},
	}

	switchCheck string
)

// ====== SearchService ======

func (s *SearchServiceMock) SearchInstruction(request m.InstructionSearchRequestDTO) (m.InstructionSearchResultDTO, error) {
	switch switchCheck {
	case "search":
		return searchResultDTO, nil
	default:
		return m.InstructionSearchResultDTO{}, errors.New("error")
	}
}

// ====== Tests ======

func TestSearch_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewSearchHandlers(&SearchServiceMock{}, &m.LoggerInterfaceMock{})

	switchCheck = "search"
	reqBody, _ := json.Marshal(searchRequestDTO)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.SearchInstruction(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(searchResultDTO)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestSearch_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewSearchHandlers(&SearchServiceMock{}, &m.LoggerInterfaceMock{})

	switchCheck = "search"

	req := httptest.NewRequest("POST", "http://example.com/api/v2/tag/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.SearchInstruction(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestSearch_SearchErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewSearchHandlers(&SearchServiceMock{}, &m.LoggerInterfaceMock{})

	switchCheck = "error"

	reqBody, _ := json.Marshal(searchRequestDTO)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.SearchInstruction(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}
