package endpoints

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_TagNotImplemented(t *testing.T) {
	e := NewTagEndpoints()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.NotImplemented(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_TagGetAll(t *testing.T) {
	e := NewTagEndpoints()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetAll(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_TagGetSingle(t *testing.T) {
	e := NewTagEndpoints()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetSingle(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_TagCreate(t *testing.T) {
	e := NewTagEndpoints()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Create(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_TagUpdate(t *testing.T) {
	e := NewTagEndpoints()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Update(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_TagDelete(t *testing.T) {
	e := NewTagEndpoints()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Delete(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}
