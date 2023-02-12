package handlers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	m "github.com/ihulsbus/cookbook/internal/models"
)

type RecipeService interface {
	FindAll() ([]m.Recipe, error)
	FindSingle(recipeID int) (m.Recipe, error)
	FindInstruction(recipeID int) ([]m.Instruction, error)
	Create(recipe m.Recipe) (m.Recipe, error)
	CreateInstruction(instruction []m.Instruction) ([]m.Instruction, error)
	Update(recipe m.Recipe) (m.Recipe, error)
	Delete(recipe m.Recipe) error
}

type ImageService interface {
	UploadImage(file multipart.File, recipeID int) bool
}

type RecipeHandlers struct {
	recipeService RecipeService
	imageService  ImageService
	logger        LoggerInterface
	utils         HanderUtils
}

func NewRecipeHandlers(recipes RecipeService, imageService ImageService, logger LoggerInterface) *RecipeHandlers {
	return &RecipeHandlers{
		recipeService: recipes,
		imageService:  imageService,
		logger:        logger,
		utils:         *NewHanderUtils(logger),
	}
}

func (h RecipeHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	var data []m.Recipe

	data, err := h.recipeService.FindAll()
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

func (h RecipeHandlers) Get(w http.ResponseWriter, r *http.Request, recipeID string) {
	var data m.Recipe

	rID, err := strconv.Atoi(recipeID)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	data, err = h.recipeService.FindSingle(rID)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

func (h RecipeHandlers) GetInstruction(w http.ResponseWriter, r *http.Request, recipeID string) {
	var data []m.Instruction

	rID, err := strconv.Atoi(recipeID)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	data, err = h.recipeService.FindInstruction(rID)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

func (h RecipeHandlers) CreateInstruction(w http.ResponseWriter, r *http.Request) {
	var instructions []m.Instruction
	var data []m.Instruction

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	if err = json.Unmarshal(body, &instructions); err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	data, err = h.recipeService.CreateInstruction(instructions)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusCreated, data)

}

func (h RecipeHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var recipe m.Recipe
	var data m.Recipe

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	if err = json.Unmarshal(body, &recipe); err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	data, err = h.recipeService.Create(recipe)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusCreated, data)
}

func (h RecipeHandlers) ImageUpload(w http.ResponseWriter, r *http.Request, recipeID string) {

	file, _, err := r.FormFile("file")
	if err != nil {
		h.utils.response400WithDetails(w, "bad request")
		return
	}

	ID, err := strconv.Atoi(recipeID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	_, err = h.recipeService.FindSingle(ID)
	if err != nil {
		h.utils.response400WithDetails(w, "recipe does not exist")
		return
	}

	if s := h.imageService.UploadImage(file, ID); s {
		h.utils.response201(w)
		return
	}

	h.utils.response500(w)

}

func (h RecipeHandlers) Update(w http.ResponseWriter, r *http.Request) {
	var recipe m.Recipe
	var data m.Recipe

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	if err = json.Unmarshal(body, &recipe); err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	if recipe.ID == 0 {
		h.utils.response400WithDetails(w, "ID is required")
		return
	}

	data, err = h.recipeService.Update(recipe)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

func (h RecipeHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	var recipe m.Recipe

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	if err = json.Unmarshal(body, &recipe); err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	if recipe.ID == 0 {
		h.utils.response400WithDetails(w, "ID is required")
		return
	}

	err = h.recipeService.Delete(recipe)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.response204(w)
}
