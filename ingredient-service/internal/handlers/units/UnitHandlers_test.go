package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	m "ingredient-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type UnitServiceMock struct {
}

var (
	units []m.UnitDTO
	unit  m.UnitDTO = m.UnitDTO{
		ID:        uuid.New(),
		FullName:  "Fluid Ounce",
		ShortName: "fl oz",
	}
)

func (s *UnitServiceMock) FindAll() ([]m.UnitDTO, error) {
	switch unit.FullName {
	case "findall":
		units = append(units, unit)
		return units, nil
	case "notfound":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (s *UnitServiceMock) FindSingle(unitDTO m.UnitDTO) (m.UnitDTO, error) {
	switch unit.FullName {
	case "find":
		return unit, nil
	case "notfound":
		return m.UnitDTO{}, errors.New("not found")
	default:
		return m.UnitDTO{}, errors.New("error")
	}
}

func (s *UnitServiceMock) Create(unitDTO m.UnitDTO) (m.UnitDTO, error) {
	switch unitDTO.FullName {
	case "create":
		return unit, nil
	default:
		return unit, errors.New("error")
	}
}

func (s *UnitServiceMock) Update(unitDTO m.UnitDTO) (m.UnitDTO, error) {
	switch unitDTO.FullName {
	case "update":
		return unit, nil
	default:
		return unit, errors.New("error")
	}
}

func (s *UnitServiceMock) Delete(unitDTO m.UnitDTO) error {
	switch unit.FullName {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

type LoggerInterfaceMock struct{}

func (l *LoggerInterfaceMock) Debugf(format string, args ...interface{}) {}
func (l *LoggerInterfaceMock) Warnf(format string, args ...interface{})  {}

// ==================================================================================================
func TestUnitGetAll_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	units = append(units, unit)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	unit.FullName = "findall"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/units", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(units)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestUnitGetAll_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	units = append(units, unit)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	unit.FullName = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/units", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no units found"}`, string(body))
}

func TestUnitGetAll_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	units = append(units, unit)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	unit.FullName = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/units", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestUnitGet_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	unit.FullName = "find"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/unit/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: unit.ID.String()},
	}

	h.GetSingle(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(unit)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestUnitGet_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	unit.FullName = "find"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/unit/", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetSingle(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid unit ID"}`, string(body))
}

func TestUnitGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	unit.FullName = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/unit/0", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: unit.ID.String()},
	}

	h.GetSingle(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no unit found"}`, string(body))
}

func TestUnitGet_FindErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	unit.FullName = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/unit/0", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: unit.ID.String()},
	}

	h.GetSingle(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestUnitCreate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	createUnit := m.UnitDTO{
		FullName: "create",
	}
	reqBody, _ := json.Marshal(createUnit)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/unit", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assertBody, _ := json.Marshal(unit)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, assertBody, body)
}

func TestUnitCreate_UnmarshallErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/unit", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unexpected JSON input"}`, string(body))
}

func TestUnitCreate_CreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	createRecipe := m.UnitDTO{
		FullName: "error",
	}
	reqBody, _ := json.Marshal(createRecipe)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/unit", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestUnitUpdate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	unit.FullName = "update"
	reqBody, _ := json.Marshal(unit)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/unit/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: unit.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestUnitUpdate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/unit/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: unit.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestUnitUpdate_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(unit)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/unit/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid unit ID"}`, string(body))
}

func TestUnitUpdate_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	unit.FullName = "fail"
	reqBody, _ := json.Marshal(unit)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/unit/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: unit.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestUnitDelete_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	unit.FullName = "delete"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/unit/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: unit.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUnitDelete_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/unit/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"invalid unit ID"}`), body)
}

func TestUnitDelete_DeleteErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUnitHandlers(&UnitServiceMock{}, &LoggerInterfaceMock{})

	unit.FullName = "error"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/unit/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: unit.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"error"}`), body)
}
