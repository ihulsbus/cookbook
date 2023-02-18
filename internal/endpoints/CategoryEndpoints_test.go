package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type CategoryHandlersMock struct{}

func (h *CategoryHandlersMock) GetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *CategoryHandlersMock) Get(w http.ResponseWriter, r *http.Request, categoryID string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *CategoryHandlersMock) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *CategoryHandlersMock) Update(w http.ResponseWriter, r *http.Request, categoryID string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
func (h *CategoryHandlersMock) Delete(w http.ResponseWriter, r *http.Request, categoryID string) {
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte(""))
}

// ==================================================================================================

func Test_CategoryGetAll(t *testing.T) {
	e := NewCategoryEndpoints(&CategoryHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetAll(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_CategoryGetSingle(t *testing.T) {
	e := NewCategoryEndpoints(&CategoryHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetSingle(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_CategoryCreate(t *testing.T) {
	e := NewCategoryEndpoints(&CategoryHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Create(c)

	assert.Equal(t, w.Code, http.StatusCreated)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_CategoryUpdate(t *testing.T) {
	e := NewCategoryEndpoints(&CategoryHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Update(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	assert.Equal(t, w.Body.String(), "{}")
}

func Test_CategoryDelete(t *testing.T) {
	e := NewCategoryEndpoints(&CategoryHandlersMock{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Delete(c)

	assert.Equal(t, w.Code, http.StatusNoContent)
}
