package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	tc "image-service/internal/test_common"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/ihulsbus/cookbook/shared/oidc"
	"github.com/stretchr/testify/assert"
)

var (
	imgs       []m.ImageDataDTO
	imgDataDTO m.ImageDataDTO
	imgFileDTO m.ImageFileDTO

	// findTrigger controls imgServiceMock.Find's dispatch independently of
	// imgDataDTO.EntityType. Defaults to "ok" so the resolve step succeeds
	// unless a test overrides it.
	findTrigger = "ok"

	testPrincipalID = uuid.New().String()
)

func init() {
	imgDataDTO = m.ImageDataDTO{
		ID:         uuid.New(),
		EntityID:   uuid.New(),
		EntityType: "recipe",
		Size:       0,
		Type:       "img/jpeg",
	}

	file, _ := tc.CreateFile()
	imgFileDTO = m.ImageFileDTO{
		ID:         imgDataDTO.ID,
		EntityID:   imgDataDTO.EntityID,
		EntityType: imgDataDTO.EntityType,
		Size:       imgDataDTO.Size,
		Type:       imgDataDTO.Type,
		File:       file,
	}
}

type imgServiceMock struct{}

type LoggerInterfaceMock struct{}

func (l *LoggerInterfaceMock) Debugf(format string, args ...interface{}) {}
func (l *LoggerInterfaceMock) Warnf(format string, args ...interface{})  {}
func (l *LoggerInterfaceMock) Errorf(format string, args ...interface{}) {}
func (l *LoggerInterfaceMock) Infof(format string, args ...interface{})  {}

// mockAuthorizer always allows by default; allow can be flipped per test
// to exercise the denial path.
type mockAuthorizer struct {
	allow bool
	err   error
}

func newAllowAllAuthorizer() *mockAuthorizer { return &mockAuthorizer{allow: true} }

func (a *mockAuthorizer) CanRead(ctx context.Context, principal *oidc.Principal, recipeID uuid.UUID) (bool, error) {
	return a.allow, a.err
}

func (a *mockAuthorizer) CanWrite(ctx context.Context, principal *oidc.Principal, recipeID uuid.UUID) (bool, error) {
	return a.allow, a.err
}

func (s *imgServiceMock) FindAll(pagination m.PaginationRequest) (m.PaginatedResponse[m.ImageDataDTO], error) {
	switch imgDataDTO.EntityType {
	case "findall":
		return m.PaginatedResponse[m.ImageDataDTO]{
			Data:       imgs,
			Pagination: m.NewPaginationMetadata(pagination.Page, pagination.Limit, int64(len(imgs))),
		}, nil
	default:
		return m.PaginatedResponse[m.ImageDataDTO]{}, errors.New("error")
	}
}

func (s *imgServiceMock) Find(imDTO m.ImageDataDTO) (m.ImageDataDTO, error) {
	switch findTrigger {
	case "ok":
		return imgDataDTO, nil
	case "notfound":
		return m.ImageDataDTO{}, errors.New("not found")
	default:
		return m.ImageDataDTO{}, errors.New("error")
	}
}

func (s *imgServiceMock) Create(imDTO m.ImageFileDTO) (m.ImageDataDTO, error) {
	switch imgDataDTO.EntityType {
	case "create":
		return imgDataDTO, nil
	default:
		return m.ImageDataDTO{}, errors.New("error")
	}
}

func (s *imgServiceMock) Update(imgDTO m.ImageFileDTO) (m.ImageDataDTO, error) {
	switch imgDataDTO.EntityType {
	case "update":
		return imgDataDTO, nil
	default:
		return m.ImageDataDTO{}, errors.New("error")
	}
}

func (s *imgServiceMock) Delete(imDTO m.ImageDataDTO) error {
	switch imgDataDTO.EntityType {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

// authedContext builds a gin context with a principal set, mimicking oidc.Module.RequireAuth.
func authedContext(w *httptest.ResponseRecorder, req *http.Request) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	oidc.SetPrincipal(c, &oidc.Principal{UserID: testPrincipalID})
	return c
}

func TestImageGetAll_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	imgs = append(imgs, imgDataDTO)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	imgDataDTO.EntityType = "findall"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.FindAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	pagination := m.NewPaginationMetadata(1, 25, int64(len(imgs)))
	expectedBody, _ := json.Marshal(m.PaginatedResponse[m.ImageDataDTO]{Data: imgs, Pagination: pagination})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestImageGetAll_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	imgs = append(imgs, imgDataDTO)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	imgDataDTO.EntityType = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.FindAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestImageGet_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	findTrigger = "ok"
	imgDataDTO.EntityType = "recipe"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDataDTO.ID.String()},
	}

	h.Find(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(imgDataDTO)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestImageGet_NotOwner_Denied(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, &mockAuthorizer{allow: false}, &LoggerInterfaceMock{})

	findTrigger = "ok"
	imgDataDTO.EntityType = "recipe"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDataDTO.ID.String()},
	}

	h.Find(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no images found"}`, string(body))
}

func TestImageGet_NonRecipeEntity_NoAuthzNeeded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Authorizer would deny, but EntityType isn't "recipe" so it's never consulted.
	h := NewHttpHandler(&imgServiceMock{}, &mockAuthorizer{allow: false}, &LoggerInterfaceMock{})

	findTrigger = "ok"
	imgDataDTO.EntityType = "ingredient"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDataDTO.ID.String()},
	}

	h.Find(c)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestImageGet_NoID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	findTrigger = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Find(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid image ID"}`, string(body))
}

func TestImageGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	findTrigger = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDataDTO.ID.String()},
	}

	h.Find(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no images found"}`, string(body))
}

func TestImageGet_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	findTrigger = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDataDTO.ID.String()},
	}

	h.Find(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestImageCreate_JpegOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	imgDataDTO.EntityType = "create"

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="image"; filename="test.jpg"`)
	header.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(header)
	writer.FormDataContentType()
	assert.NoError(t, err)
	part.Write(buf.Bytes())
	writer.Close()

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDataDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDataDTO.EntityType},
	}

	h.Create(c)

	resp := w.Result()
	expectedBody, err := json.Marshal(imgDataDTO)
	assert.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestImageCreate_PngOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	imgDataDTO.EntityType = "create"

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	assert.NoError(t, err)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="image"; filename="test.png"`)
	header.Set("Content-Type", "image/png")
	part, err := writer.CreatePart(header)
	writer.FormDataContentType()
	assert.NoError(t, err)
	part.Write(buf.Bytes())
	writer.Close()

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDataDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDataDTO.EntityType},
	}

	h.Create(c)

	resp := w.Result()
	expectedBody, err := json.Marshal(imgDataDTO)
	assert.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestImageCreate_NotOwner_Denied(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, &mockAuthorizer{allow: false}, &LoggerInterfaceMock{})

	entityID := uuid.New()

	req := httptest.NewRequest("POST", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: entityID.String()},
		gin.Param{Key: "entityType", Value: "recipe"},
	}

	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no images found"}`, string(body))
}

func TestImageCreate_EntityIDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="imag"; filename="test.jpg"`)
	header.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(header)
	writer.FormDataContentType()
	assert.NoError(t, err)
	part.Write(buf.Bytes())
	writer.Close()

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	imgDataDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid entityID"}`, string(body))
}

func TestImageCreate_EntityTypeErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="imag"; filename="test.jpg"`)
	header.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(header)
	writer.FormDataContentType()
	assert.NoError(t, err)
	part.Write(buf.Bytes())
	writer.Close()

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDataDTO.ID.String()},
	}

	imgDataDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"entityType is required"}`, string(body))
}

func TestImageCreate_NoFileErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="imag"; filename="test.jpg"`)
	header.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(header)
	writer.FormDataContentType()
	assert.NoError(t, err)
	part.Write(buf.Bytes())
	writer.Close()

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDataDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDataDTO.EntityType},
	}

	imgDataDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"image file is required"}`, string(body))
}

func TestImageCreate_ImageFileTooLargeErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	img := image.NewRGBA(image.Rect(0, 0, 19000, 19000))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="image"; filename="test.jpg"`)
	header.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(header)
	writer.FormDataContentType()
	assert.NoError(t, err)
	part.Write(buf.Bytes())
	writer.Close()

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDataDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDataDTO.EntityType},
	}

	imgDataDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"image file is too large"}`, string(body))
}

func TestImageCreate_InvalidFormatErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="image"; filename="test.jpg"`)
	header.Set("Content-Type", "image/invalid")
	part, err := writer.CreatePart(header)
	writer.FormDataContentType()
	assert.NoError(t, err)
	part.Write(buf.Bytes())
	writer.Close()

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDataDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDataDTO.EntityType},
	}

	imgDataDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unsupported image format: image/invalid"}`, string(body))
}

func TestImageCreate_DecodeErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="image"; filename="test.jpg"`)
	header.Set("Content-Type", "image/png")
	part, err := writer.CreatePart(header)
	writer.FormDataContentType()
	assert.NoError(t, err)
	part.Write(buf.Bytes())
	writer.Close()

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDataDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDataDTO.EntityType},
	}

	imgDataDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid image: png: invalid format: not a PNG file"}`, string(body))
}

func TestImageCreate_ImageTooSmallErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	img := image.NewRGBA(image.Rect(0, 0, 299, 299))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="image"; filename="test.jpg"`)
	header.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(header)
	writer.FormDataContentType()
	assert.NoError(t, err)
	part.Write(buf.Bytes())
	writer.Close()

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDataDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDataDTO.EntityType},
	}

	imgDataDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"image dimensions are too small. Dimensions need to be between 300x300 and 1000x1000"}`, string(body))
}

func TestImageCreate_ImageTooBigErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	img := image.NewRGBA(image.Rect(0, 0, 1001, 1001))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="image"; filename="test.jpg"`)
	header.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(header)
	writer.FormDataContentType()
	assert.NoError(t, err)
	part.Write(buf.Bytes())
	writer.Close()

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDataDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDataDTO.EntityType},
	}

	imgDataDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"image dimensions are too big. Dimensions need to be between 300x300 and 1000x1000"}`, string(body))
}

func TestImageCreate_ImageCreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	img := image.NewRGBA(image.Rect(0, 0, 300, 300))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="image"; filename="test.jpg"`)
	header.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(header)
	writer.FormDataContentType()
	assert.NoError(t, err)
	part.Write(buf.Bytes())
	writer.Close()

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDataDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDataDTO.EntityType},
	}

	imgDataDTO.EntityType = "error"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestImageUpdate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	findTrigger = "ok"

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="image"; filename="test.jpg"`)
	header.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(header)
	writer.FormDataContentType()
	assert.NoError(t, err)
	part.Write(buf.Bytes())
	writer.Close()

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDataDTO.ID.String()},
	}

	imgDataDTO.EntityType = "update"
	h.Update(c)

	resp := w.Result()
	expBody, _ := json.Marshal(imgDataDTO)
	fmt.Println(string(expBody))
	respBody, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expBody, respBody)
}

func TestImageUpdate_NotOwner_Denied(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, &mockAuthorizer{allow: false}, &LoggerInterfaceMock{})

	findTrigger = "ok"
	imgDataDTO.EntityType = "recipe"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", bytes.NewReader(nil))
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDataDTO.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no images found"}`, string(body))
}

func TestImageUpdate_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	imgDataDTO.EntityType = "update"
	reqBody, _ := json.Marshal(imgDataDTO)

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid image ID"}`, string(body))
}

func TestImageUpdate_ResolveErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	findTrigger = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDataDTO.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no images found"}`, string(body))
}

func TestImageUpdate_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	findTrigger = "ok"

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="image"; filename="test.jpg"`)
	header.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(header)
	writer.FormDataContentType()
	assert.NoError(t, err)
	part.Write(buf.Bytes())
	writer.Close()

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDataDTO.ID.String()},
	}

	imgDataDTO.EntityType = "fail"
	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestImageDelete_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	findTrigger = "ok"
	imgDataDTO.EntityType = "delete"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDataDTO.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestImageDelete_NotOwner_Denied(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, &mockAuthorizer{allow: false}, &LoggerInterfaceMock{})

	findTrigger = "ok"
	imgDataDTO.EntityType = "recipe"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDataDTO.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no images found"}`, string(body))
}

func TestImageDelete_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	imgDataDTO.EntityType = "delete"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"invalid image ID"}`), body)
}

func TestImageDelete_ResolveErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	findTrigger = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDataDTO.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"no images found"}`), body)
}

func TestImageDelete_DeleteErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHttpHandler(&imgServiceMock{}, newAllowAllAuthorizer(), &LoggerInterfaceMock{})

	findTrigger = "ok"
	imgDataDTO.EntityType = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c := authedContext(w, req)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDataDTO.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"error"}`), body)
}
