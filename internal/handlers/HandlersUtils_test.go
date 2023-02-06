package handlers

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse201(t *testing.T) {
	h := NewHanderUtils(&LoggerInterfaceMock{})
	w := httptest.NewRecorder()

	h.response201(w)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, body, []byte(`"Object created"`))
}

func TestResponse204(t *testing.T) {
	h := NewHanderUtils(&LoggerInterfaceMock{})
	w := httptest.NewRecorder()

	h.response204(w)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, w.Code, 204)
	assert.Equal(t, body, []byte(`"Object deleted"`))
}

func TestResponse500(t *testing.T) {
	h := NewHanderUtils(&LoggerInterfaceMock{})
	w := httptest.NewRecorder()

	h.response500(w)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, w.Code, 500)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal server error"}`))
}

func TestResponse400WithDetails(t *testing.T) {
	h := NewHanderUtils(&LoggerInterfaceMock{})
	w := httptest.NewRecorder()

	h.response400WithDetails(w, "400")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (400)"}`))
}

func TestResponse500WithDetails(t *testing.T) {
	h := NewHanderUtils(&LoggerInterfaceMock{})
	w := httptest.NewRecorder()

	h.response500WithDetails(w, "500")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, w.Code, 500)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (500)"}`))
}

func TestRespondWithError(t *testing.T) {
	h := NewHanderUtils(&LoggerInterfaceMock{})
	w := httptest.NewRecorder()

	h.respondWithError(w, 500, "500")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, w.Code, 500)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"500"}`))
}

func TestRespondWithJSON(t *testing.T) {
	h := NewHanderUtils(&LoggerInterfaceMock{})
	w := httptest.NewRecorder()

	h.respondWithJSON(w, 200, "200")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, body, []byte(`"200"`))
}
