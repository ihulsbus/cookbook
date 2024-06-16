package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	m "metadata-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TagServiceMock struct {
}

var (
	tags []m.TagDTO
	tag  m.TagDTO = m.TagDTO{
		ID:   uuid.New(),
		Name: "tag",
	}
)

// ====== TagService ======

func (s *TagServiceMock) FindAll() ([]m.TagDTO, error) {
	switch tag.Name {
	case "findall":
		return tags, nil
	case "notfound":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (s *TagServiceMock) FindSingle(tagDTO m.TagDTO) (m.TagDTO, error) {
	switch tag.Name {
	case "find":
		return tag, nil
	case "notfound":
		return m.TagDTO{}, errors.New("not found")
	default:
		return m.TagDTO{}, errors.New("error")
	}
}

func (s *TagServiceMock) Create(tagDTO m.TagDTO) (m.TagDTO, error) {
	switch tagDTO.Name {
	case "create":
		return tag, nil
	default:
		return m.TagDTO{}, errors.New("error")
	}
}

func (s *TagServiceMock) Update(tagDTO m.TagDTO) (m.TagDTO, error) {
	switch tagDTO.Name {
	case "update":
		return tag, nil
	default:
		return m.TagDTO{}, errors.New("error")
	}
}

func (s *TagServiceMock) Delete(tagDTO m.TagDTO) error {
	switch tag.Name {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

func TestTagGetAll_OK(t *testing.T) {
	tags = append(tags, tag)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	tag.Name = "findall"

	req := httptest.NewRequest("GET", "http://example.com/api/v1/tag", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(tags)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestTagGetAll_NotFound(t *testing.T) {
	tags = append(tags, tag)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	tag.Name = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v1/tag", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no tags found"}`, string(body))
}

func TestTagGetAll_Error(t *testing.T) {
	tags = append(tags, tag)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	tag.Name = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v1/tag", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestTagGet_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: tag.ID.String()},
	}

	tag.Name = "find"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(tag)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestTagGet_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	tag.Name = "find"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid tag ID"}`, string(body))
}

func TestTagGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: tag.ID.String()},
	}

	tag.Name = "notfound"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"tag not found"}`, string(body))
}

func TestTagGet_FindErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: tag.ID.String()},
	}

	tag.Name = "finderr"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestTagCreate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	createTag := m.TagDTO{
		Name: "create",
	}
	reqBody, _ := json.Marshal(createTag)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assertBody, _ := json.Marshal(tag)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, assertBody, body)
}

func TestTagCreate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v1/tag/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unexpected JSON input"}`, string(body))
}

func TestTagCreate_CreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	createTag := m.TagDTO{
		Name: "createerr",
	}
	reqBody, _ := json.Marshal(createTag)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestTagUpdate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	tag.Name = "update"
	reqBody, _ := json.Marshal(tag)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: tag.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestTagUpdate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/tag/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: tag.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestTagUpdate_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(tag)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid tag ID"}`, string(body))
}

func TestTagUpdate_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	tag.Name = "updatefail"
	reqBody, _ := json.Marshal(tag)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: tag.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestTagDelete_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	tag.Name = "delete"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: tag.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, []byte(``), body)
}

func TestTagDelete_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"invalid tag ID"}`), body)
}

func TestTagDelete_DeleteErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	tag.Name = "deleteError"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: tag.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"error"}`), body)
}
