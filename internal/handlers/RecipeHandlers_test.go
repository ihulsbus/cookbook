package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
		StepNumber:  1,
		Description: "instruction",
	}
)

// ====== RecipeService ======

func (s *RecipeServiceMock) FindAll() ([]m.Recipe, error) {
	return recipes, nil
}

func (s *RecipeServiceMock) FindSingle(recipeID int) (m.Recipe, error) {
	switch recipeID {
	case 1:
		return m.Recipe{RecipeName: "recipe1"}, nil
	case 2:
		return m.Recipe{RecipeName: "recipe2"}, nil
	default:
		return m.Recipe{}, errors.New("error")
	}
}

func (s *RecipeServiceMock) FindInstruction(recipeID int) ([]m.Instruction, error) {
	var instructions []m.Instruction
	switch recipeID {
	case 1:
		instructions := append(instructions, m.Instruction{})
		return instructions, nil
	case 2:
		instructions := append(instructions, m.Instruction{})
		return instructions, nil
	default:
		return nil, errors.New("error")
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

func (s *RecipeServiceMock) CreateInstruction(instructions []m.Instruction) ([]m.Instruction, error) {
	switch instructions[0].RecipeID {
	case 1:
		return instructions, nil
	default:
		return nil, errors.New("error")
	}
}

func (s *RecipeServiceMock) Update(recipe m.Recipe) (m.Recipe, error) {
	switch recipe.RecipeName {
	case "recipe":
		return recipe, nil
	default:
		return recipe, errors.New("error")
	}
}

func (s *RecipeServiceMock) Delete(recipe m.Recipe) error {
	switch recipe.ID {
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
	assert.Equal(t, string(body), `{"code":500,"msg":"Internal Server Error. (error)"}`)
}

func TestRecipeGetInstruction_OK(t *testing.T) {
	var instructions []m.Instruction
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.GetInstruction(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	instructions = append(instructions, m.Instruction{})
	expectedBody, _ := json.Marshal(instructions)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBody, body)
}

func TestRecipeGetInstruction_AtoiErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("GET", "http://example.com/api/v1/recipe/1", nil)
	w := httptest.NewRecorder()

	h.GetInstruction(w, req, "")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody := `{"code":500,"msg":"Internal Server Error. (strconv.Atoi: parsing \"\": invalid syntax)"}`

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, expectedBody, string(body))
}

func TestRecipeGetInstruction_FindErr(t *testing.T) {
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

func TestRecipeCreate_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	reqBody, _ := json.Marshal(recipe)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, body, reqBody)
}

func TestRecipeCreate_UnmarshalErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Create(w, req)

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

	h.Create(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestRecipeCreateInstruction_OK(t *testing.T) {
	var createInstruction []m.Instruction
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	createInstruction = append(createInstruction, instruction)

	reqBody, _ := json.Marshal(createInstruction)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1/instruction", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.CreateInstruction(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, body, reqBody)
}

func TestRecipeCreateInstruction_UnmarshalErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.CreateInstruction(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (unexpected end of JSON input)"}`))
}

func TestRecipeCreateInstruction_CreateErr(t *testing.T) {
	var createInstruction []m.Instruction
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	cI := instruction
	cI.RecipeID = 2
	createInstruction = append(createInstruction, cI)

	reqBody, _ := json.Marshal(createInstruction)

	req := httptest.NewRequest("POST", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.CreateInstruction(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	fmt.Print(string(body))

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
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

	h.ImageUpload(w, req, "0")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (recipe does not exist)"}`))
}

func TestRecipeUpdate_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	updateRecipe := recipe
	updateRecipe.ID = 1
	updateRecipe.RecipeName = "recipe"
	reqBody, _ := json.Marshal(updateRecipe)

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Update(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, reqBody)
}

func TestRecipeUpdate_UnmarshalErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("PUT", "http://example.com/api/v1/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Update(w, req)

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

	h.Update(w, req)

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

	h.Update(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
}

func TestRecipeDelete_OK(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	deleteRecipe := recipe
	deleteRecipe.ID = 1
	deleteRecipe.RecipeName = "recipe"
	reqBody, _ := json.Marshal(deleteRecipe)

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)
}

func TestRecipeDelete_UnmarshalErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/recipe/1", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (unexpected end of JSON input)"}`))
}

func TestRecipeDelete_IDRequiredErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	deleteRecipe := recipe
	deleteRecipe.ID = 0
	deleteRecipe.RecipeName = "recipe"
	reqBody, _ := json.Marshal(deleteRecipe)

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	assert.Equal(t, body, []byte(`{"code":400,"msg":"Bad Request. (ID is required)"}`))
}

func TestRecipeDelete_DeleteErr(t *testing.T) {
	h := NewRecipeHandlers(&RecipeServiceMock{}, &ImageServiceMock{}, &LoggerInterfaceMock{})

	deleteRecipe := recipe
	deleteRecipe.ID = 2
	deleteRecipe.RecipeName = "recipe"
	reqBody, _ := json.Marshal(deleteRecipe)

	req := httptest.NewRequest("DELETE", "http://example.com/api/v1/recipe/1", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, body, []byte(`{"code":500,"msg":"Internal Server Error. (error)"}`))
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
