package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
)

var (
	instruction = models.InstructionDTO{
		ID:          uuid.New(),
		Sequence:    1,
		Description: "instruction",
		MediaID:     uuid.New(),
		EntityID:    uuid.New(),
		EntityType:  "recipe",
	}
)

type InstructionServiceMock struct {
}

func (s *InstructionServiceMock) Find(recipeID uuid.UUID) (*[]models.InstructionDTO, error) {
	switch instruction.Description {
	case "find":
		var instructions []models.InstructionDTO
		instructions = append(instructions, instruction)
		return &instructions, nil
	case "notfound":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (s *InstructionServiceMock) Create(entityID uuid.UUID, instructionDTO *[]models.InstructionDTO) (*[]models.InstructionDTO, error) {
	switch instruction.Description {
	case "create":
		var instructions []models.InstructionDTO
		instructions = append(instructions, instruction)
		return &instructions, nil
	default:
		return nil, errors.New("error")
	}
}

func (s *InstructionServiceMock) Update(entityID uuid.UUID, instructionDTO *[]models.InstructionDTO) (*[]models.InstructionDTO, error) {
	switch instruction.Description {
	case "update":
		var instructions []models.InstructionDTO
		instructions = append(instructions, instruction)
		return &instructions, nil
	default:
		return nil, errors.New("error")
	}
}

func (s *InstructionServiceMock) Delete(recipeID uuid.UUID) error {
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
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

	instruction.Description = "find"

	req := httptest.NewRequest("GET", fmt.Sprintf("http://example.com/api/v2/instruction/%s", instruction.EntityID), nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: instruction.ID.String()},
	}

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal([]models.InstructionDTO{instruction})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestGetInstruction_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

	instruction.Description = "create"
	reqBody, _ := json.Marshal([]models.InstructionDTO{instruction})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/instruction/1/instruction", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.AddParam("id", "6c4e174f-e760-4a4d-af6e-d1d5849b0fe1")

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assertBody, _ := json.Marshal([]models.InstructionDTO{instruction})

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, assertBody, body)
}

func TestCreateInstruction_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/instruction/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.AddParam("id", "6c4e174f-e760-4a4d-af6e-d1d5849b0fe1")

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unexpected JSON input"}`, string(body))
}

func TestCreateInstruction_CreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

	instruction.Description = "createError"
	reqBody, _ := json.Marshal([]models.InstructionDTO{instruction})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/instruction/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.AddParam("id", "6c4e174f-e760-4a4d-af6e-d1d5849b0fe1")

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestUpdateInstruction_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

	instruction.Description = "update"
	reqBody, _ := json.Marshal([]models.InstructionDTO{instruction})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.AddParam("id", "6c4e174f-e760-4a4d-af6e-d1d5849b0fe1")

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestUpdateEntity_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

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
	assert.Equal(t, `{"error":"invalid entity id provided"}`, string(body))
}

func TestUpdateInstruction_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.AddParam("recipeID", "6c4e174f-e760-4a4d-af6e-d1d5849b0fe1")

	h.Update(c)

	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateInstruction_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

	instruction.Description = "updateError"
	reqBody, _ := json.Marshal([]models.InstructionDTO{instruction})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/instruction/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.AddParam("id", "6c4e174f-e760-4a4d-af6e-d1d5849b0fe1")

	h.Update(c)

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteInstruction_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

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
	h := NewInstructionHandlers(&InstructionServiceMock{}, &models.LoggerInterfaceMock{})

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
