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

var (
	instruction m.InstructionDTO = m.InstructionDTO{
		ID:          uuid.New(),
		Sequence:    1,
		Description: "instruction",
		MediaID:     uuid.New(),
	}
)

type InstructionServiceMock struct {
}

func (s *InstructionServiceMock) Find(instructionDTO m.InstructionDTO) (m.InstructionDTO, error) {
	switch instruction.Description {
	case "find":
		return instruction, nil
	case "notfound":
		return m.InstructionDTO{}, errors.New("not found")
	default:
		return m.InstructionDTO{}, errors.New("error")
	}
}

func (s *InstructionServiceMock) Create(instructionDTO m.InstructionDTO) (m.InstructionDTO, error) {
	switch instructionDTO.Description {
	case "create":
		return instruction, nil
	default:
		return m.InstructionDTO{}, errors.New("error")
	}
}

func (s *InstructionServiceMock) Update(instructionDTO m.InstructionDTO) (m.InstructionDTO, error) {
	switch instructionDTO.Description {
	case "update":
		return instruction, nil
	default:
		return m.InstructionDTO{}, errors.New("error")
	}
}

func (s *InstructionServiceMock) Delete(instructionDTO m.InstructionDTO) error {
	switch instruction.Description {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

// ========================================================================================================

func TestGetInstruction_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	instruction.Description = "find"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: instruction.ID.String()},
	}

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(instruction)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestGetInstruction_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	instruction.Description = "find"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid instruction ID"}`, string(body))
}

func TestGetInstruction_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	instruction.Description = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: instruction.ID.String()},
	}

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no instruction found"}`, string(body))
}

func TestGetInstruction_FindErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	instruction.Description = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: instruction.ID.String()},
	}

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestCreateInstruction_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	createInstruction := m.InstructionDTO{
		Sequence:    1,
		Description: "create",
		MediaID:     instruction.MediaID,
	}
	reqBody, _ := json.Marshal(createInstruction)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/instruction/1/instruction", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assertBody, _ := json.Marshal(instruction)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, assertBody, body)
}

func TestCreateInstruction_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/instruction/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unexpected JSON input"}`, string(body))
}

func TestCreateInstruction_CreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	createInstruction := m.InstructionDTO{
		Sequence:    1,
		Description: "error",
		MediaID:     instruction.MediaID,
	}
	reqBody, _ := json.Marshal(createInstruction)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/instruction/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestUpdateInstruction_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	instruction.Description = "update"
	reqBody, _ := json.Marshal(instruction)

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: instruction.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestUpdateInstruction_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	instruction.Description = "update"
	reqBody, _ := json.Marshal(instruction)

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid instruction ID"}`, string(body))
}

func TestUpdateInstruction_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: instruction.ID.String()},
	}

	h.Update(c)

	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateInstruction_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	instruction.Description = "error"
	reqBody, _ := json.Marshal(instruction)

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: instruction.ID.String()},
	}

	h.Update(c)

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteInstruction_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	instruction.Description = "delete"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: instruction.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteInstruction_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	instruction.Description = "delete"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid instruction ID"}`, string(body))
}

func TestDeleteInstruction_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &m.LoggerInterfaceMock{})

	instruction.Description = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: instruction.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}
