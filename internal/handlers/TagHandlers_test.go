package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	m "github.com/ihulsbus/cookbook/internal/models"
	"github.com/stretchr/testify/assert"
)

type TagServiceMock struct {
}

var (
	tags []m.Tag
	tag  m.Tag = m.Tag{
		TagName: "tag",
	}
)

// ====== TagService ======

func (s *TagServiceMock) FindAll() ([]m.Tag, error) {
	return tags, nil
}

func (s *TagServiceMock) FindSingle(tagID uint) (m.Tag, error) {
	switch tagID {
	case 1:
		return m.Tag{TagName: "tag1"}, nil
	case 2:
		return m.Tag{TagName: "tag2"}, nil
	default:
		return m.Tag{}, errors.New("error")
	}
}

func (s *TagServiceMock) Create(tag m.Tag) (m.Tag, error) {
	switch tag.TagName {
	case "tag":
		return tag, nil
	default:
		return tag, errors.New("error")
	}
}

func (s *TagServiceMock) Update(tag m.Tag, tagID uint) (m.Tag, error) {
	switch tag.TagName {
	case "tag":
		return tag, nil
	default:
		return tag, errors.New("error")
	}
}

func (s *TagServiceMock) Delete(tagID uint) error {
	switch tagID {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

func TestTagGetAll_OK(t *testing.T) {
	tags = append(tags, m.Tag{TagName: "tag1"}, m.Tag{TagName: "tag2"})
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/tag", nil)
	w := httptest.NewRecorder()

	h.GetAll(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(tags)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestTagGet_OK(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()

	h.Get(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(m.Tag{TagName: "tag1"})

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestTagGet_AtoiErr(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()

	h.Get(w, req, "")

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
}

func TestTagGet_FindErr(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()

	h.Get(w, req, "3")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, `{"code":500,"msg":"Internal Server Error. (error)"}`, string(body))
}

func TestTagCreate_OK(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(tag)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, body, reqBody)
}

func TestTagCreate_UnmarshalErr(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v1/tag/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (unexpected end of JSON input)"}`))
}

func TestTagCreate_CreateErr(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	tag.TagName = "err"
	reqBody, _ := json.Marshal(tag)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestTagUpdate_OK(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	updateTag := tag
	updateTag.ID = 1
	updateTag.TagName = "tag"
	reqBody, _ := json.Marshal(updateTag)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, reqBody)
}

func TestTagUpdate_AtoiErr(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/tag/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Update(w, req, "")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, []byte(`{"code":500,"msg":"Internal server error"}`), body)
}

func TestTagUpdate_UnmarshalErr(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/tag/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (unexpected end of JSON input)"}`))
}

func TestTagUpdate_IDRequiredErr(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	updateTag := tag
	updateTag.ID = 0
	updateTag.TagName = "tag"
	reqBody, _ := json.Marshal(updateTag)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestTagUpdate_UpdateErr(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	updateTag := tag
	updateTag.ID = 2
	updateTag.TagName = "fail"
	reqBody, _ := json.Marshal(updateTag)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/tag/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestTagDelete_OK(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "1")

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)
}

func TestTagDelete_AtoiErr(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "")

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestTagDelete_IDRequiredErr(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestTagDelete_DeleteErr(t *testing.T) {
	h := NewTagHandlers(&TagServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/tag/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "2")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}
