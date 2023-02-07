package endpoints

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_CategoryNotImplemented(t *testing.T) {
	e := NewCategoryEndpoints()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.NotImplemented(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_CategoryGetAll(t *testing.T) {
	e := NewCategoryEndpoints()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetAll(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_CategoryGetSingle(t *testing.T) {
	e := NewCategoryEndpoints()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.GetSingle(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_CategoryCreate(t *testing.T) {
	e := NewCategoryEndpoints()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Create(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_CategoryUpdate(t *testing.T) {
	e := NewCategoryEndpoints()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Update(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}

func Test_CategoryDelete(t *testing.T) {
	e := NewCategoryEndpoints()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e.Delete(c)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, `"not implemented"`, w.Body.String())
}
