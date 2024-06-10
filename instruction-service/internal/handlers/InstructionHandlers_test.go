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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	instruction m.Instruction = m.Instruction{
		RecipeID:    uuid.UUID{},
		Description: "instruction",
	}
)

type InstructionServiceMock struct {
}

func (s *InstructionServiceMock) FindInstruction(recipeID uuid.UUID) (m.Instruction, error) {
	switch recipeID {
	case uuid.UUID{}:
		return m.Instruction{}, nil
	case uuid.UUID{}:
		return m.Instruction{}, nil
	default:
		return m.Instruction{}, errors.New("error")
	}
}

func (s *InstructionServiceMock) CreateInstruction(instruction m.Instruction, recipeID uuid.UUID) (m.Instruction, error) {
	switch instruction.RecipeID {
	case uuid.UUID{}:
		return instruction, nil
	default:
		return m.Instruction{}, errors.New("error")
	}
}

func (s *InstructionServiceMock) UpdateInstruction(instruction m.Instruction, recipeID uuid.UUID) (m.Instruction, error) {
	switch recipeID {
	case uuid.UUID{}:
		return instruction, nil
	default:
		return m.Instruction{}, errors.New("error")
	}
}

func (s *InstructionServiceMock) DeleteInstruction(recipeID uuid.UUID) error {
	switch recipeID {
	case uuid.UUID{}:
		return nil
	default:
		return errors.New("error")
	}
}

// ========================================================================================================

func TestGetInstruction_OK(t *testing.T) {
	h := NewInstructionHandlers(&InstructionServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.GetInstruction(w, req, uuid.UUID{})

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(m.Instruction{})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestGetInstruction_FindErr(t *testing.T) {
	h := NewInstructionHandlers(&InstructionServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.GetInstruction(w, req, uuid.UUID{})

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody := `{"code":500,"msg":"Internal Server Error. (error)"}`

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, expectedBody, string(body))
}

func TestCreateInstruction_OK(t *testing.T) {
	h := NewInstructionHandlers(&InstructionServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(instruction)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1/instruction", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.CreateInstruction(w, req, uuid.UUID{})

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, body, reqBody)
}

func TestCreateInstruction_UnmarshalErr(t *testing.T) {
	h := NewInstructionHandlers(&InstructionServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.CreateInstruction(w, req, uuid.UUID{})

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (unexpected end of JSON input)"}`))
}

func TestCreateInstruction_CreateErr(t *testing.T) {
	h := NewInstructionHandlers(&InstructionServiceMock{}, &LoggerInterfaceMock{})

	cI := instruction
	cI.RecipeID = uuid.UUID{}

	reqBody, _ := json.Marshal(cI)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.CreateInstruction(w, req, uuid.UUID{})

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestUpdateInstruction_OK(t *testing.T) {
	h := NewInstructionHandlers(&InstructionServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(instruction)

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.UpdateInstruction(w, req, uuid.UUID{})

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestUpdateInstruction_UnmarshalErr(t *testing.T) {
	h := NewInstructionHandlers(&InstructionServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.UpdateInstruction(w, req, uuid.UUID{})

	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateInstruction_UpdateErr(t *testing.T) {
	h := NewInstructionHandlers(&InstructionServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(instruction)

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.UpdateInstruction(w, req, uuid.UUID{})

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteInstruction_OK(t *testing.T) {
	h := NewInstructionHandlers(&InstructionServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.DeleteInstruction(w, req, uuid.UUID{})

	resp := w.Result()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestDeleteInstruction_UpdateErr(t *testing.T) {
	h := NewInstructionHandlers(&InstructionServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.DeleteInstruction(w, req, uuid.UUID{})

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
