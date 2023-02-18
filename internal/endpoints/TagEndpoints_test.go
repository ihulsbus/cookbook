package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TagHandlersMock struct{}

func (h *TagHandlersMock) GetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *TagHandlersMock) Get(w http.ResponseWriter, r *http.Request, tagID string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *TagHandlersMock) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *TagHandlersMock) Update(w http.ResponseWriter, r *http.Request, tagID string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *TagHandlersMock) Delete(w http.ResponseWriter, r *http.Request, tagID string) {
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte(""))
}

// ==================================================================================================

func Test_TagGetAll(t *testing.T) {
	e := NewTagEndpoints(&TagHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetAll(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_TagGetSingle(t *testing.T) {
	e := NewTagEndpoints(&TagHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetSingle(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_TagCreate(t *testing.T) {
	e := NewTagEndpoints(&TagHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Create(c)

	assert.Equal(t, w.Code, http.StatusCreated)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_TagUpdate(t *testing.T) {
	e := NewTagEndpoints(&TagHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Update(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_TagDelete(t *testing.T) {
	e := NewTagEndpoints(&TagHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Delete(c)

	assert.Equal(t, w.Code, http.StatusNoContent)
}
