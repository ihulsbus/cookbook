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
	"github.com/stretchr/testify/assert"
)

type MetadataServiceMock struct {
}

var (
	metadata = models.RecipeMetadataDTO{
		RecipeID: uuid.New(),
		CuisineType: models.CuisineTypeDTO{
			ID:   uuid.New(),
			Name: "Cuisine",
		},
		DifficultyLevel: models.DifficultyLevelDTO{
			ID:    uuid.New(),
			Level: 1,
		},
		Tags:            []models.TagDTO{{ID: uuid.New(), Name: "Tag"}},
		PreparationTime: 99,
		ServingCount:    99,
	}
)

// ====== MetadataService ======

func (s *MetadataServiceMock) FindAll(pagination models.PaginationRequest) (models.PaginatedResponse[models.RecipeMetadataDTO], error) {
	switch metadata.ServingCount {
	case 1: // find all
		metadataArray := []models.RecipeMetadataDTO{metadata}
		return models.PaginatedResponse[models.RecipeMetadataDTO]{
			Data:       metadataArray,
			Pagination: models.NewPaginationMetadata(pagination.Page, pagination.Limit, int64(len(metadataArray))),
		}, nil
	default:
		return models.PaginatedResponse[models.RecipeMetadataDTO]{}, errors.New("error")
	}
}

func (s *MetadataServiceMock) Find(recipeID uuid.UUID) (*models.RecipeMetadataDTO, error) {
	switch metadata.ServingCount {
	case 2: // find
		return &metadata, nil
	case 0: // not found
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (s *MetadataServiceMock) Create(recipeID uuid.UUID, meta *models.RecipeMetadataDTO) (*models.RecipeMetadataDTO, error) {
	switch metadata.ServingCount {
	case 3: // create
		return &metadata, nil
	default:
		return nil, errors.New("error")
	}
}

func (s *MetadataServiceMock) Update(recipeID uuid.UUID, meta *models.RecipeMetadataDTO) (*models.RecipeMetadataDTO, error) {
	switch metadata.ServingCount {
	case 4: // update
		return &metadata, nil
	default:
		return nil, errors.New("error")
	}
}

func (s *MetadataServiceMock) Delete(recipeID uuid.UUID) error {
	switch metadata.ServingCount {
	case 5: // delete
		return nil
	default:
		return errors.New("error")
	}
}

func TestMetadataGetAll_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	metadata.ServingCount = 1
	metadataArray := []models.RecipeMetadataDTO{metadata}

	req := httptest.NewRequest("GET", "http://example.com/api/v2/metadata", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	pagination := models.NewPaginationMetadata(1, 25, int64(len(metadataArray)))
	expectedResponse := models.PaginatedResponse[models.RecipeMetadataDTO]{Data: metadataArray, Pagination: pagination}
	expectedBody, _ := json.Marshal(expectedResponse)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestMetadataGetAll_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	metadata.ServingCount = 99

	req := httptest.NewRequest("GET", "http://example.com/api/v2/metadata", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestMetadataGet_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	metadata.ServingCount = 2

	req := httptest.NewRequest("GET", "http://example.com/api/v2/metadata/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: metadata.RecipeID.String()},
	}

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(metadata)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestMetadataGet_ID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/metadata/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid Recipe ID"}`, string(body))
}

func TestMetadataGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	metadata.ServingCount = 0
	req := httptest.NewRequest("GET", "http://example.com/api/v2/metadata/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: metadata.RecipeID.String()},
	}

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"recipe not found"}`, string(body))
}

func TestMetadataGet_FindErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	metadata.ServingCount = 99
	req := httptest.NewRequest("GET", "http://example.com/api/v2/metadata/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: metadata.RecipeID.String()},
	}

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestMetadataCreate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	metadata.ServingCount = 3
	reqBody, _ := json.Marshal(metadata)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/metadata/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: metadata.RecipeID.String()},
	}

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assertBody, _ := json.Marshal(metadata)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, assertBody, body)
}

func TestMetadataCreate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/metadata/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: metadata.RecipeID.String()},
	}

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestMetadataCreate_CreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	metadata.ServingCount = 99
	reqBody, _ := json.Marshal(metadata)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/metadata/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: metadata.RecipeID.String()},
	}

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestMetadataUpdate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	metadata.ServingCount = 4
	reqBody, _ := json.Marshal(metadata)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/metadata/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: metadata.RecipeID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestMetadataUpdate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/metadata/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: metadata.RecipeID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestMetadataUpdate_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(metadata)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/metadata/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid recipe ID"}`, string(body))
}

func TestMetadataUpdate_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	metadata.ServingCount = 99
	reqBody, _ := json.Marshal(metadata)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/metadata/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: metadata.RecipeID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestMetadataDelete_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	metadata.ServingCount = 5

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/metadata/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: metadata.RecipeID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, ``, string(body))
}

func TestMetadataDelete_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/metadata/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid recipe ID"}`, string(body))
}

func TestMetadataDelete_DeleteErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMetadataHandlers(&MetadataServiceMock{}, &models.LoggerInterfaceMock{})

	metadata.ServingCount = 99

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/metadata/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: metadata.RecipeID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}
