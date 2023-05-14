package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	m "github.com/ihulsbus/cookbook/internal/models"
	"github.com/stretchr/testify/assert"
)

type RecipeServiceMock struct {
}

type ImageServiceMock struct {
}

var (
	recipes []m.Recipe
	recipe  m.Recipe = m.Recipe{
		RecipeName: "recipe",
	}
	instruction m.Instruction = m.Instruction{
		RecipeID:    1,
		Description: "instruction",
	}
)

// ====== RecipeService ======

func (s *RecipeServiceMock) FindAll() ([]m.Recipe, error) {
	return recipes, nil
}

func (s *RecipeServiceMock) FindSingle(recipeID uint) (m.Recipe, error) {
	switch recipeID {
	case 1:
		return m.Recipe{RecipeName: "recipe1"}, nil
	case 2:
		return m.Recipe{RecipeName: "recipe2"}, nil
	default:
		return m.Recipe{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) Create(recipe m.Recipe) (m.Recipe, error) {
	switch recipe.RecipeName {
	case "recipe":
		return recipe, nil
	default:
		return recipe, errors.New("error")
	}
}

func (s *RecipeServiceMock) Update(recipe m.Recipe, recipeID uint) (m.Recipe, error) {
	switch recipe.RecipeName {
	case "recipe":
		return recipe, nil
	default:
		return recipe, errors.New("error")
	}
}

func (s *RecipeServiceMock) Delete(recipeID uint) error {
	switch recipeID {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

func (s *RecipeServiceMock) FindInstruction(recipeID uint) (m.Instruction, error) {
	switch recipeID {
	case 1:
		return m.Instruction{}, nil
	case 2:
		return m.Instruction{}, nil
	default:
		return m.Instruction{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) CreateInstruction(instruction m.Instruction, recipeID uint) (m.Instruction, error) {
	switch instruction.RecipeID {
	case 1:
		return instruction, nil
	default:
		return m.Instruction{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) UpdateInstruction(instruction m.Instruction, recipeID uint) (m.Instruction, error) {
	switch recipeID {
	case 1:
		return instruction, nil
	default:
		return m.Instruction{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) DeleteInstruction(recipeID uint) error {
	switch recipeID {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

func (s *RecipeServiceMock) FindIngredientLink(recipeID uint) ([]m.RecipeIngredient, error) {
	switch recipeID {
	case 1:
		return []m.RecipeIngredient{}, nil
	default:
		return nil, errors.New("error")
	}
}

func (s *RecipeServiceMock) CreateIngredientLink(link []m.RecipeIngredient, recipeID uint) ([]m.RecipeIngredient, error) {
	switch recipeID {
	case 1:
		return []m.RecipeIngredient{}, nil
	default:
		return nil, errors.New("error")
	}
}

func (s *RecipeServiceMock) UpdateIngredientLink(link []m.RecipeIngredient, recipeID uint) ([]m.RecipeIngredient, error) {
	switch recipeID {
	case 1:
		return []m.RecipeIngredient{}, nil
	default:
		return nil, errors.New("error")
	}
}

func (s *RecipeServiceMock) DeleteIngredientLink(link []m.RecipeIngredient, recipeID uint) error {
	switch recipeID {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

// ====== ImageService ======

func (S *ImageServiceMock) UploadImage(file multipart.File, recipeID int) bool {
	switch recipeID {
	case 1:
		return true
	default:
		return false
	}
}

// ==================================================================================================

func TestRecipeGetAll_OK(t *testing.T) {
	recipes = append(recipes, m.Recipe{RecipeName: "recipe1"}, m.Recipe{RecipeName: "recipe2"})
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe", nil)
	w := httptest.NewRecorder()

	h.GetAll(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(recipes)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestRecipeGet_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Get(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(m.Recipe{RecipeName: "recipe1"})

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestRecipeGet_AtoiErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Get(w, req, "")

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
}

func TestRecipeGet_FindErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Get(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, `{"code":500,"msg":"Internal Server Error. (error)"}`, string(body))
}

func TestRecipeGetIngredients_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.GetIngredientLink(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal([]m.RecipeIngredient{})

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, expectedBody)
}

func TestRecipeGetIngredients_AtoiErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.GetIngredientLink(w, req, "")

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
}

func TestRecipeGetIngredients_FindErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.GetIngredientLink(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, `{"code":500,"msg":"Internal Server Error. (error)"}`, string(body))
}

func TestRecipeCreate_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Create(&m.User{}, w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, body, reqBody)
}

func TestRecipeCreate_UnmarshalErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Create(&m.User{}, w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (unexpected end of JSON input)"}`))
}

func TestRecipeCreate_CreateErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	recipe.RecipeName = "err"
	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Create(&m.User{}, w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestRecipeImageUpload_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	// Set up a pipe to avoid buffering
	pr, pw := io.Pipe()
	// This writer is going to transform
	// what we pass to it to multipart form data
	// and write it to our io.Pipe
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		// We create the form data field 'fileupload'
		// which returns another writer to write the actual file
		part, err := writer.CreateFormFile("file", "someimg.png")
		if err != nil {
			t.Error(err)
		}

		// https://yourbasic.org/golang/create-image/
		img := createImage()

		// Encode() takes an io.Writer.
		// We pass the multipart field
		// 'fileupload' that we defined
		// earlier which, in turn, writes
		// to our io.Pipe
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1/upload", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.ImageUpload(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, body, []byte(`"Object created"`))
}

func TestRecipeImageUpload_FormErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	// Set up a pipe to avoid buffering
	pr, pw := io.Pipe()
	// This writer is going to transform
	// what we pass to it to multipart form data
	// and write it to our io.Pipe
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		// We create the form data field 'fileupload'
		// which returns another writer to write the actual file
		part, err := writer.CreateFormFile("file", "someimg.png")
		if err != nil {
			t.Error(err)
		}

		// https://yourbasic.org/golang/create-image/
		img := createImage()

		// Encode() takes an io.Writer.
		// We pass the multipart field
		// 'fileupload' that we defined
		// earlier which, in turn, writes
		// to our io.Pipe
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1/upload", pr)
	w := httptest.NewRecorder()

	h.ImageUpload(w, req, "1")

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestRecipeImageUpload_AtoiErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	// Set up a pipe to avoid buffering
	pr, pw := io.Pipe()
	// This writer is going to transform
	// what we pass to it to multipart form data
	// and write it to our io.Pipe
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		// We create the form data field 'fileupload'
		// which returns another writer to write the actual file
		part, err := writer.CreateFormFile("file", "someimg.png")
		if err != nil {
			t.Error(err)
		}

		// https://yourbasic.org/golang/create-image/
		img := createImage()

		// Encode() takes an io.Writer.
		// We pass the multipart field
		// 'fileupload' that we defined
		// earlier which, in turn, writes
		// to our io.Pipe
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe//upload", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.ImageUpload(w, req, "")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal server error"}`))
}

func TestRecipeImageUpload_RecipeDoesNotExistErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	// Set up a pipe to avoid buffering
	pr, pw := io.Pipe()
	// This writer is going to transform
	// what we pass to it to multipart form data
	// and write it to our io.Pipe
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		// We create the form data field 'fileupload'
		// which returns another writer to write the actual file
		part, err := writer.CreateFormFile("file", "someimg.png")
		if err != nil {
			t.Error(err)
		}

		// https://yourbasic.org/golang/create-image/
		img := createImage()

		// Encode() takes an io.Writer.
		// We pass the multipart field
		// 'fileupload' that we defined
		// earlier which, in turn, writes
		// to our io.Pipe
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/0/upload", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.ImageUpload(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (recipe does not exist)"}`))
}

func TestRecipeImageUpload_UploadErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	// Set up a pipe to avoid buffering
	pr, pw := io.Pipe()
	// This writer is going to transform
	// what we pass to it to multipart form data
	// and write it to our io.Pipe
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		// We create the form data field 'fileupload'
		// which returns another writer to write the actual file
		part, err := writer.CreateFormFile("file", "someimg.png")
		if err != nil {
			t.Error(err)
		}

		// https://yourbasic.org/golang/create-image/
		img := createImage()

		// Encode() takes an io.Writer.
		// We pass the multipart field
		// 'fileupload' that we defined
		// earlier which, in turn, writes
		// to our io.Pipe
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/2/upload", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.ImageUpload(w, req, "2")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, []byte(`{"code":500,"msg":"Internal server error"}`), body)
}

func TestRecipeUpdate_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	updateRecipe := recipe
	updateRecipe.ID = 1
	updateRecipe.RecipeName = "recipe"
	reqBody, _ := json.Marshal(updateRecipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, reqBody)
}

func TestRecipeUpdate_UnmarshalErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (unexpected end of JSON input)"}`))
}

func TestRecipeUpdate_IDRequiredErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	updateRecipe := recipe
	updateRecipe.ID = 0
	updateRecipe.RecipeName = "recipe"
	reqBody, _ := json.Marshal(updateRecipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestRecipeUpdate_UpdateErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	updateRecipe := recipe
	updateRecipe.ID = 2
	updateRecipe.RecipeName = "fail"
	reqBody, _ := json.Marshal(updateRecipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestRecipeDelete_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "1")

	resp := w.Result()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestRecipeDelete_AtoiErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "")

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestRecipeDelete_IDRequiredErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestRecipeDelete_DeleteErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req, "2")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

// ========================================================================================================

func TestGetInstruction_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.GetInstruction(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody, _ := json.Marshal(m.Instruction{})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestGetInstruction_AtoiErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.GetInstruction(w, req, "")

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetInstruction_FindErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.GetInstruction(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody := `{"code":500,"msg":"Internal Server Error. (error)"}`

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, expectedBody, string(body))
}

func TestCreateInstruction_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(instruction)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1/instruction", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.CreateInstruction(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, body, reqBody)
}

func TestCreateInstruction_UnmarshalErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.CreateInstruction(w, req, "2")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (unexpected end of JSON input)"}`))
}

func TestCreateInstruction_CreateErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	cI := instruction
	cI.RecipeID = 2

	reqBody, _ := json.Marshal(cI)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.CreateInstruction(w, req, "2")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestUpdateInstruction_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(instruction)

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.UpdateInstruction(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, reqBody, body)
}

func TestUpdateInstruction_AtoiErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.UpdateInstruction(w, req, "")

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestUpdateInstruction_UnmarshalErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.UpdateInstruction(w, req, "1")

	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateInstruction_UpdateErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(instruction)

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.UpdateInstruction(w, req, "2")

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteInstruction_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.DeleteInstruction(w, req, "1")

	resp := w.Result()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestDeleteInstruction_AtoiErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.DeleteInstruction(w, req, "")

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteInstruction_UpdateErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.DeleteInstruction(w, req, "2")

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// ====== Helpers ======
func createImage() *image.RGBA {
	width := 200
	height := 100

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{100, 200, 200, 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}

	return img
}
