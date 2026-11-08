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
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/ihulsbus/cookbook/shared/oidc"
	"github.com/stretchr/testify/assert"
)

type RecipeServiceMock struct {
}

var (
	testRequesterID uuid.UUID = uuid.New()

	recipes []m.RecipeDTO
	recipe  m.RecipeDTO = m.RecipeDTO{
		ID:           uuid.New(),
		Name:         "recipe",
		Description:  "description",
		ServingCount: 1,
		OwnerID:      testRequesterID,
		Visibility:   m.VisibilityPrivate,
	}
)

// authedContext builds a gin test context with testRequesterID set as the
// authenticated principal, as oidc.Module.RequireAuth would before any
// handler in this package runs.
func authedContext(w *httptest.ResponseRecorder, req *http.Request) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	oidc.SetPrincipal(c, &oidc.Principal{UserID: testRequesterID.String()})
	return c
}

// ====== RecipeService ======

func (s *RecipeServiceMock) FindAll(pagination m.PaginationRequest, requesterID uuid.UUID, requesterGroups []string, scope m.RecipeListScope) (m.PaginatedResponse[m.RecipeDTO], error) {
	switch recipe.Name {
	case "findall":
		return m.PaginatedResponse[m.RecipeDTO]{
			Data:       recipes,
			Pagination: m.NewPaginationMetadata(pagination.Page, pagination.Limit, int64(len(recipes))),
		}, nil
	case "invalidscope":
		return m.PaginatedResponse[m.RecipeDTO]{}, errors.New("invalid scope")
	default:
		return m.PaginatedResponse[m.RecipeDTO]{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) FindSingle(recipeDTO m.RecipeDTO, requesterID uuid.UUID, requesterGroups []string) (m.RecipeDTO, error) {
	switch recipe.Name {
	case "find":
		return recipe, nil
	case "notfound":
		return m.RecipeDTO{}, errors.New("not found")
	default:
		return m.RecipeDTO{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) Create(recipeDTO m.RecipeDTO, ownerID uuid.UUID) (m.RecipeDTO, error) {
	switch recipeDTO.Name {
	case "create":
		return recipe, nil
	default:
		return m.RecipeDTO{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) Update(recipeDTO m.RecipeDTO, requesterID uuid.UUID, requesterGroups []string) (m.RecipeDTO, error) {
	switch recipeDTO.Name {
	case "update":
		return recipe, nil
	case "notfound":
		return m.RecipeDTO{}, errors.New("not found")
	default:
		return m.RecipeDTO{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) Delete(recipeDTO m.RecipeDTO, requesterID uuid.UUID) error {
	switch recipe.Name {
	case "delete":
		return nil
	case "notfound":
		return errors.New("not found")
	default:
		return errors.New("error")
	}
}

// shareTrigger controls ListShares/UpsertShare/DeleteShare/UpdateVisibility
// mock responses - set per test before calling the handler, since these
// handlers' only input besides the body is the recipe ID path param
// (shareTrigger stands in for a dedicated dispatch field, mirroring how
// the other mocks above dispatch on recipe.Name/recipeDTO.Name).
var shareTrigger string

func (s *RecipeServiceMock) ListShares(recipeID uuid.UUID, requesterID uuid.UUID) ([]m.RecipeShareDTO, error) {
	switch shareTrigger {
	case "ok":
		return []m.RecipeShareDTO{{TargetType: m.ShareTargetUser, TargetID: uuid.New().String(), Permission: m.SharePermissionRead}}, nil
	case "notfound":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (s *RecipeServiceMock) UpsertShare(recipeID uuid.UUID, requesterID uuid.UUID, share m.RecipeShareDTO) (m.RecipeShareDTO, error) {
	switch shareTrigger {
	case "ok":
		return share, nil
	case "notfound":
		return m.RecipeShareDTO{}, errors.New("not found")
	case "invalid target type", "invalid permission", "target id is required", "target id must be a valid user id":
		return m.RecipeShareDTO{}, errors.New(shareTrigger)
	default:
		return m.RecipeShareDTO{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) DeleteShare(recipeID uuid.UUID, requesterID uuid.UUID, targetType m.ShareTargetType, targetID string) error {
	switch shareTrigger {
	case "ok":
		return nil
	case "notfound":
		return errors.New("not found")
	default:
		return errors.New("error")
	}
}

func (s *RecipeServiceMock) UpdateVisibility(recipeID uuid.UUID, requesterID uuid.UUID, visibility m.Visibility) (m.RecipeDTO, error) {
	switch shareTrigger {
	case "ok":
		updated := recipe
		updated.Visibility = visibility
		return updated, nil
	case "notfound":
		return m.RecipeDTO{}, errors.New("not found")
	case "invalid visibility":
		return m.RecipeDTO{}, errors.New("invalid visibility")
	default:
		return m.RecipeDTO{}, errors.New("error")
	}
}

// Logging
type LoggerInterfaceMock struct{}

func (l *LoggerInterfaceMock) Debugf(format string, args ...interface{}) {}
func (l *LoggerInterfaceMock) Warnf(format string, args ...interface{})  {}
func (l *LoggerInterfaceMock) Errorf(format string, args ...interface{}) {}
func (l *LoggerInterfaceMock) Infof(format string, args ...interface{})  {}

// ==================================================================================================

func TestPrincipalID_MissingPrincipal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w) // no principal set
	c.Request = req

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, `{"error":"missing authenticated principal"}`, string(body))
}

func TestPrincipalID_NonUUIDSub(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	oidc.SetPrincipal(c, &oidc.Principal{UserID: "not-a-uuid"})

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid principal"}`, string(body))
}

func TestRecipeGetAll_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recipes = append(recipes, recipe)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "findall"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	pagination := m.NewPaginationMetadata(1, 25, int64(len(recipes)))
	expectedBody, _ := json.Marshal(m.PaginatedResponse[m.RecipeDTO]{Data: recipes, Pagination: pagination})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestRecipeGetAll_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recipes = append(recipes, recipe)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

// TestRecipeGetAll_InvalidScope: the service's "invalid scope" error
// (KIT-30) maps to 400, not the generic 500 every other service error gets.
func TestRecipeGetAll_InvalidScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "invalidscope"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe?scope=bogus", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.GetAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid scope"}`, string(body))
}

func TestRecipeGet_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	recipe.Name = "find"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(recipe)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestRecipeGet_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	recipe.Name = "find"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid recipe ID"}`, string(body))
}

func TestRecipeGet_notFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	recipe.Name = "notfound"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"recipe not found"}`, string(body))
}

func TestRecipeGet_FindErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	recipe.Name = "finderr"

	h.Get(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestRecipeCreate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	createRecipe := m.RecipeDTO{
		Name: "create",
	}
	reqBody, _ := json.Marshal(createRecipe)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assertBody, _ := json.Marshal(recipe)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, assertBody, body)
}

func TestRecipeCreate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v2/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unexpected JSON input"}`, string(body))
}

func TestRecipeCreate_CreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	createRecipe := m.RecipeDTO{
		Name: "error",
	}
	reqBody, _ := json.Marshal(createRecipe)

	req := httptest.NewRequest("POST", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestRecipeUpdate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "update"
	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestRecipeUpdate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestRecipeUpdate_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid recipe ID"}`, string(body))
}

func TestRecipeUpdate_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "fail"
	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestRecipeUpdate_NotFoundErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "notfound"
	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"recipe not found"}`, string(body))
}

func TestRecipeDelete_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "delete"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRecipeDelete_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"invalid recipe ID"}`), body)
}

func TestRecipeDelete_DeleteErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "error"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"error"}`), body)
}

func TestRecipeDelete_NotFoundErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	recipe.Name = "notfound"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"recipe not found"}`), body)
}

func TestRecipeGetShares_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "ok"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1/shares", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{gin.Param{Key: "id", Value: recipe.ID.String()}}

	h.GetShares(c)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRecipeGetShares_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1/shares", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.GetShares(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid recipe ID"}`, string(body))
}

func TestRecipeGetShares_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1/shares", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{gin.Param{Key: "id", Value: recipe.ID.String()}}

	h.GetShares(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"recipe not found"}`, string(body))
}

func TestRecipeGetShares_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/recipe/1/shares", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{gin.Param{Key: "id", Value: recipe.ID.String()}}

	h.GetShares(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestRecipeUpsertShare_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "ok"

	shareDTO := m.RecipeShareDTO{TargetType: m.ShareTargetUser, TargetID: uuid.New().String(), Permission: m.SharePermissionRead}
	reqBody, _ := json.Marshal(shareDTO)

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1/shares", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{gin.Param{Key: "id", Value: recipe.ID.String()}}

	h.UpsertShare(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	expectedBody, _ := json.Marshal(shareDTO)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestRecipeUpsertShare_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1/shares", bytes.NewReader([]byte(`{}`)))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.UpsertShare(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid recipe ID"}`, string(body))
}

func TestRecipeUpsertShare_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1/shares", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{gin.Param{Key: "id", Value: recipe.ID.String()}}

	h.UpsertShare(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unexpected JSON input"}`, string(body))
}

func TestRecipeUpsertShare_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "notfound"

	reqBody, _ := json.Marshal(m.RecipeShareDTO{TargetType: m.ShareTargetUser, TargetID: uuid.New().String(), Permission: m.SharePermissionRead})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1/shares", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{gin.Param{Key: "id", Value: recipe.ID.String()}}

	h.UpsertShare(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"recipe not found"}`, string(body))
}

func TestRecipeUpsertShare_ValidationErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "invalid permission"

	reqBody, _ := json.Marshal(m.RecipeShareDTO{TargetType: m.ShareTargetUser, TargetID: uuid.New().String(), Permission: "nonsense"})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1/shares", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{gin.Param{Key: "id", Value: recipe.ID.String()}}

	h.UpsertShare(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid permission"}`, string(body))
}

func TestRecipeUpsertShare_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "error"

	reqBody, _ := json.Marshal(m.RecipeShareDTO{TargetType: m.ShareTargetUser, TargetID: uuid.New().String(), Permission: m.SharePermissionRead})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1/shares", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{gin.Param{Key: "id", Value: recipe.ID.String()}}

	h.UpsertShare(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestRecipeDeleteShare_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "ok"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1/shares/user/2", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
		gin.Param{Key: "targetType", Value: "user"},
		gin.Param{Key: "targetId", Value: uuid.New().String()},
	}

	h.DeleteShare(c)

	resp := w.Result()

	// Matches TestRecipeDelete_OK's assertion for the same reason: this
	// harness's httptest.ResponseRecorder only picks up the real status
	// code once something is actually written to the body, and
	// ctx.Status(204) with no body never triggers that - it's seen here
	// as 200, not the real 204 DeleteShare sends.
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRecipeDeleteShare_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1/shares/user/2", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.DeleteShare(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid recipe ID"}`, string(body))
}

func TestRecipeDeleteShare_InvalidTargetType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1/shares/bogus/2", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
		gin.Param{Key: "targetType", Value: "bogus"},
		gin.Param{Key: "targetId", Value: uuid.New().String()},
	}

	h.DeleteShare(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid target type"}`, string(body))
}

func TestRecipeDeleteShare_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "notfound"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1/shares/user/2", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
		gin.Param{Key: "targetType", Value: "user"},
		gin.Param{Key: "targetId", Value: uuid.New().String()},
	}

	h.DeleteShare(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"recipe not found"}`, string(body))
}

func TestRecipeDeleteShare_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "error"

	req := httptest.NewRequest("DELETE", "http://example.com/api/v2/recipe/1/shares/user/2", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: recipe.ID.String()},
		gin.Param{Key: "targetType", Value: "user"},
		gin.Param{Key: "targetId", Value: uuid.New().String()},
	}

	h.DeleteShare(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestRecipeUpdateVisibility_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "ok"

	reqBody, _ := json.Marshal(map[string]string{"visibility": "public"})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1/visibility", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{gin.Param{Key: "id", Value: recipe.ID.String()}}

	h.UpdateVisibility(c)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRecipeUpdateVisibility_IDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1/visibility", bytes.NewReader([]byte(`{}`)))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.UpdateVisibility(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid recipe ID"}`, string(body))
}

func TestRecipeUpdateVisibility_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1/visibility", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{gin.Param{Key: "id", Value: recipe.ID.String()}}

	h.UpdateVisibility(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unexpected JSON input"}`, string(body))
}

func TestRecipeUpdateVisibility_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "notfound"

	reqBody, _ := json.Marshal(map[string]string{"visibility": "public"})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1/visibility", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{gin.Param{Key: "id", Value: recipe.ID.String()}}

	h.UpdateVisibility(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"recipe not found"}`, string(body))
}

func TestRecipeUpdateVisibility_InvalidVisibility(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "invalid visibility"

	reqBody, _ := json.Marshal(map[string]string{"visibility": "nonsense"})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1/visibility", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{gin.Param{Key: "id", Value: recipe.ID.String()}}

	h.UpdateVisibility(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid visibility"}`, string(body))
}

func TestRecipeUpdateVisibility_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandlers(&RecipeServiceMock{}, &LoggerInterfaceMock{})

	shareTrigger = "error"

	reqBody, _ := json.Marshal(map[string]string{"visibility": "public"})

	req := httptest.NewRequest("PUT", "http://example.com/api/v2/recipe/1/visibility", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{gin.Param{Key: "id", Value: recipe.ID.String()}}

	h.UpdateVisibility(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}
