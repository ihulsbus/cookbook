package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type InstructionHandlersMock struct {
}

type MiddlewareMock struct{}

func (h *InstructionHandlersMock) GetInstruction(w http.ResponseWriter, r *http.Request, recipeID uuid.UUID) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *InstructionHandlersMock) CreateInstruction(w http.ResponseWriter, r *http.Request, recipeID uuid.UUID) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *InstructionHandlersMock) UpdateInstruction(w http.ResponseWriter, r *http.Request, recipeID uuid.UUID) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *InstructionHandlersMock) DeleteInstruction(w http.ResponseWriter, r *http.Request, recipeID uuid.UUID) {
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte(""))
}

// ==================================================================================================

func Test_GetInstruction(t *testing.T) {
	e := NewInstructionEndpoints(&InstructionHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetInstruction(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_CreateInstruction(t *testing.T) {
	e := NewInstructionEndpoints(&InstructionHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.CreateInstruction(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_UpdateInstruction(t *testing.T) {
	e := NewInstructionEndpoints(&InstructionHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.UpdateInstruction(c)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_DeleteInstruction(t *testing.T) {
	e := NewInstructionEndpoints(&InstructionHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.DeleteInstruction(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
