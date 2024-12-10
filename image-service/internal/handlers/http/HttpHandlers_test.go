package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	m "image-service/internal/models"
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
	"github.com/stretchr/testify/assert"
)

var (
	imgs   []m.ImageDTO
	imgDTO m.ImageDTO = m.ImageDTO{
		ID:         uuid.New(),
		EntityType: "recipe",
		EntityID:   uuid.New(),
		Size:       0,
		Type:       "img/jpeg",
		File:       tc.CreateFile(),
	}
)

type imgServiceMock struct{}

type LoggerInterfaceMock struct{}

func (l *LoggerInterfaceMock) Debugf(format string, args ...interface{}) {}
func (l *LoggerInterfaceMock) Warnf(format string, args ...interface{})  {}
func (l *LoggerInterfaceMock) Errorf(format string, args ...interface{}) {}
func (l *LoggerInterfaceMock) Infof(format string, args ...interface{})  {}

func (s *imgServiceMock) FindAll() ([]m.ImageDTO, error) {
	switch imgDTO.EntityType {
	case "findall":
		return imgs, nil
	case "notfound":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (s *imgServiceMock) Find(imDTO m.ImageDTO) (m.ImageDTO, error) {
	switch imgDTO.EntityType {
	case "find":
		return imgDTO, nil
	case "notfound":
		return m.ImageDTO{}, errors.New("not found")
	default:
		return m.ImageDTO{}, errors.New("error")
	}
}

func (s *imgServiceMock) Create(imDTO m.ImageDTO) (m.ImageDTO, error) {
	switch imgDTO.EntityType {
	case "create":
		return imgDTO, nil
	default:
		return m.ImageDTO{}, errors.New("error")
	}
}

func (s *imgServiceMock) Update(imgDTO m.ImageDTO) (m.ImageDTO, error) {
	switch imgDTO.EntityType {
	case "update":
		return imgDTO, nil
	default:
		return m.ImageDTO{}, errors.New("error")
	}
}

func (s *imgServiceMock) Delete(imDTO m.ImageDTO) error {
	switch imgDTO.EntityType {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

// ==================================================================================================

func TestImageGetAll_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	imgs = append(imgs, imgDTO)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "findall"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.FindAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(imgs)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestImageGetAll_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	imgs = append(imgs, imgDTO)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.FindAll(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no images found"}`, string(body))
}

func TestImageGetAll_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	imgs = append(imgs, imgDTO)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "error"

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
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "find"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDTO.ID.String()},
	}

	h.Find(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(imgDTO)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestImageGet_NoID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Find(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid image ID"}`, string(body))
}

func TestImageGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "notfound"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDTO.ID.String()},
	}

	h.Find(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, `{"error":"no images found"}`, string(body))
}

func TestImageGet_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDTO.ID.String()},
	}

	h.Find(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestImageCreate_JpegOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "create"

	// Create a sample image
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	// Create a new multipart writer
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDTO.EntityType},
	}

	h.Create(c)

	resp := w.Result()
	expectedBody, err := json.Marshal(imgDTO)
	assert.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestImageCreate_PngOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "create"

	// Create a sample image
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	assert.NoError(t, err)

	// Create a new multipart writer
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDTO.EntityType},
	}

	h.Create(c)

	resp := w.Result()
	expectedBody, err := json.Marshal(imgDTO)
	assert.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestImageCreate_EntityIDErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	// Create a sample image
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	// Create a new multipart writer
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	imgDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid entityID"}`, string(body))
}

func TestImageCreate_EntityTypeErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	// Create a sample image
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	// Create a new multipart writer
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDTO.ID.String()},
	}

	imgDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"entityType is required"}`, string(body))
}

func TestImageCreate_NoFileErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	// Create a sample image
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	// Create a new multipart writer
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDTO.EntityType},
	}

	imgDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"image file is required"}`, string(body))
}

func TestImageCreate_ImageFileTooLargeErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	// Create a sample image
	img := image.NewRGBA(image.Rect(0, 0, 19000, 19000))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	// Create a new multipart writer
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDTO.EntityType},
	}

	imgDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"image file is too large"}`, string(body))
}

func TestImageCreate_InvalidFormatErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	// Create a sample image
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	// Create a new multipart writer
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDTO.EntityType},
	}

	imgDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"unsupported image format: image/invalid"}`, string(body))
}

func TestImageCreate_DecodeErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	// Create a sample image
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	// Create a new multipart writer
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDTO.EntityType},
	}

	imgDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid image"}`, string(body))
}

func TestImageCreate_ImageTooSmallErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	// Create a sample image
	img := image.NewRGBA(image.Rect(0, 0, 299, 299))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	// Create a new multipart writer
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDTO.EntityType},
	}

	imgDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"image dimensions are too small"}`, string(body))
}

func TestImageCreate_ImageTooBigErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	// Create a sample image
	img := image.NewRGBA(image.Rect(0, 0, 1001, 1001))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	// Create a new multipart writer
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDTO.EntityType},
	}

	imgDTO.EntityType = "create"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"image dimensions are too big"}`, string(body))
}

func TestImageCreate_ImageCreateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	// Create a sample image
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	assert.NoError(t, err)

	// Create a new multipart writer
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "entityID", Value: imgDTO.ID.String()},
		gin.Param{Key: "entityType", Value: imgDTO.EntityType},
	}

	imgDTO.EntityType = "error"
	h.Create(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestImageUpdate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "update"
	reqBody, _ := json.Marshal(imgDTO)

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDTO.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestImageUpdate_UnmarshalErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDTO.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"EOF"}`, string(body))
}

func TestImageUpdate_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "update"
	reqBody, _ := json.Marshal(imgDTO)

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"invalid image ID"}`, string(body))
}

func TestImageUpdate_UpdateErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "fail"
	reqBody, _ := json.Marshal(imgDTO)

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDTO.ID.String()},
	}

	h.Update(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error":"error"}`, string(body))
}

func TestImageDelete_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "delete"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDTO.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestImageDelete_IDRequiredErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "delete"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"invalid image ID"}`), body)
}

func TestImageDelete_DeleteErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewImageHandlers(&imgServiceMock{}, &LoggerInterfaceMock{})

	imgDTO.EntityType = "error"

	req := httptest.NewRequest("GET", "http://example.com/api/v2/img", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: imgDTO.ID.String()},
	}

	h.Delete(c)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"error":"error"}`), body)
}
