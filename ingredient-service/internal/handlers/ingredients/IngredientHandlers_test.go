package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
	"github.com/ihulsbus/cookbook/shared/oidc"
	"github.com/stretchr/testify/assert"
)

var testCreatedBy = uuid.New()

// authedContext builds a gin test context with testCreatedBy set as the
// authenticated principal, as oidc.Module.RequireAuth would before this runs.
func authedContext(w *httptest.ResponseRecorder, req *http.Request) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	oidc.SetPrincipal(c, &oidc.Principal{UserID: testCreatedBy.String()})
	return c
}

type IngredientServiceMock struct {
}

var (
	ingredients []models.IngredientDTO
	ingredient  models.IngredientDTO = models.IngredientDTO{
		ID:   uuid.New(),
		Name: "ingredient",
	}

	units []models.UnitDTO
	unit  models.UnitDTO = models.UnitDTO{
		ID:        uuid.New(),
		FullName:  "Fluid Ounce",
		ShortName: "fl oz",
	}
)

func (s *IngredientServiceMock) FindAll(pagination models.PaginationRequest) (models.PaginatedResponse[models.IngredientDTO], error) {
	switch ingredient.Name {
	case "findall":
		ingredients = append(ingredients, ingredient)
		return models.PaginatedResponse[models.IngredientDTO]{
			Data:       ingredients,
			Pagination: models.NewPaginationMetadata(pagination.Page, pagination.Limit, int64(len(ingredients))),
		}, nil
	default:
		return models.PaginatedResponse[models.IngredientDTO]{}, errors.New("error")
	}
}

func (s *IngredientServiceMock) FindSingle(ingredientDTO models.IngredientDTO) (models.IngredientDTO, error) {
	switch ingredient.Name {
	case "find":
		return ingredient, nil
	case "notfound":
		return models.IngredientDTO{}, errors.New("not found")
	default:
		return models.IngredientDTO{}, errors.New("error")
	}
}

func (s *IngredientServiceMock) Create(ingredientDTO models.IngredientDTO, createdBy uuid.UUID, force bool) (models.IngredientDTO, []models.IngredientDTO, error) {
	switch ingredientDTO.Name {
	case "create":
		return ingredient, nil, nil
	case "duplicate":
		return models.IngredientDTO{}, []models.IngredientDTO{ingredient}, errors.New("duplicate")
	case "similar":
		if force {
			return ingredient, nil, nil
		}
		return models.IngredientDTO{}, []models.IngredientDTO{ingredient}, errors.New("similar ingredients exist")
	case "":
		return models.IngredientDTO{}, nil, errors.New("name is empty")
	default:
		return ingredient, nil, errors.New("error")
	}
}

func (s *IngredientServiceMock) Suggest(q string) ([]models.IngredientDTO, error) {
	switch q {
	case "":
		return nil, errors.New("query is required")
	case "error":
		return nil, errors.New("error")
	default:
		return []models.IngredientDTO{ingredient}, nil
	}
}

func (s *IngredientServiceMock) Update(ingredientDTO models.IngredientDTO) (models.IngredientDTO, error) {
	switch ingredientDTO.Name {
	case "update":
		return ingredient, nil
	default:
		return ingredient, errors.New("error")
	}
}

func (s *IngredientServiceMock) Delete(ingredientDTO models.IngredientDTO) error {
	switch ingredient.Name {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

type LoggerInterfaceMock struct{}

func (l *LoggerInterfaceMock) Debugf(format string, args ...interface{}) {}
func (l *LoggerInterfaceMock) Warnf(format string, args ...interface{})  {}
func (l *LoggerInterfaceMock) Infof(format string, args ...interface{})  {}
func (l *LoggerInterfaceMock) Errorf(format string, args ...interface{}) {}

func TestIngredientGetAll_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ingredients = nil
	ingredients = append(ingredients, ingredient)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	ingredient.Name = "findall"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredients", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	pagination := models.NewPaginationMetadata(1, 25, int64(len(ingredients)))
	expectedBody, _ := json.Marshal(models.PaginatedResponse[models.IngredientDTO]{Data: ingredients, Pagination: pagination})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestIngredientGetAll_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	ingredient.Name = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredients", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestIngredientGet_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	ingredient.Name = "find"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredient/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: ingredient.ID.String()},
	}

	h.GetSingle(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(ingredient)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestIngredientGet_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	ingredient.Name = "find"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredient/", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetSingle(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid ingredient ID"}`, string(body))
}

func TestIngredientGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	ingredient.Name = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredient/0", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: ingredient.ID.String()},
	}

	h.GetSingle(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no ingredient found"}`, string(body))
}

func TestIngredientGet_FindErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	ingredient.Name = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredient/0", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: ingredient.ID.String()},
	}

	h.GetSingle(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestIngredientCreate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	createIngredient := models.IngredientDTO{
		Name: "create",
	}
	reqBody, _ := json.Marshal(createIngredient)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/ingredient", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assertBody, _ := json.Marshal(ingredient)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, assertBody, body)
}

func TestIngredientCreate_NoPrincipal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/ingredient", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestIngredientCreate_UnmarshallErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/ingredient", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unexpected JSON input"}`, string(body))
}

func TestIngredientCreate_CreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	createRecipe := models.IngredientDTO{
		Name: "error",
	}
	reqBody, _ := json.Marshal(createRecipe)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/ingredient", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

// An exact normalized-name match: 409, with the existing ingredient as a candidate.
func TestIngredientCreate_Duplicate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(models.IngredientDTO{Name: "duplicate"})
	req := httptest.NewRequest("POST", "http://example.com/api/v2/ingredient", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Create(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"ingredient already exists"`)
}

// A near-duplicate without force=true: 409, with candidates to pick from.
func TestIngredientCreate_Similar(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(models.IngredientDTO{Name: "similar"})
	req := httptest.NewRequest("POST", "http://example.com/api/v2/ingredient?force=false", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Create(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "did you mean")
}

// The same near-duplicate, but force=true bypasses the block.
func TestIngredientCreate_SimilarForced(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(models.IngredientDTO{Name: "similar"})
	req := httptest.NewRequest("POST", "http://example.com/api/v2/ingredient?force=true", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIngredientSuggest_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredient/suggest?q=tomat", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Suggest(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIngredientSuggest_EmptyQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/ingredient/suggest", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Suggest(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIngredientUpdate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	ingredient.Name = "update"
	reqBody, _ := json.Marshal(ingredient)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/ingredient/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: ingredient.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestIngredientUpdate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/ingredient/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: ingredient.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestIngredientUpdate_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(ingredient)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/ingredient/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid ingredient ID"}`, string(body))
}

func TestIngredientUpdate_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	ingredient.Name = "fail"
	reqBody, _ := json.Marshal(ingredient)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/ingredient/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: ingredient.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestIngredientDelete_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	ingredient.Name = "delete"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/ingredient/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: ingredient.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestIngredientDelete_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/ingredient/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"invalid ingredient ID"}`), body)
}

func TestIngredientDelete_DeleteErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewIngredientHandlers(&IngredientServiceMock{}, &LoggerInterfaceMock{})

	ingredient.Name = "error"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/ingredient/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: ingredient.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"error"}`), body)
}
