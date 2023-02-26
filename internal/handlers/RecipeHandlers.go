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
	FindSingle(recipeID uint) (m.Recipe, error)
	FindRecipeIngredients(recipeID uint) ([]m.RecipeIngredient, error)
	Create(recipe m.Recipe) (m.Recipe, error)
	Update(recipe m.Recipe, recipeID uint) (m.Recipe, error)
	Delete(recipeID uint) error

	FindInstruction(recipeID uint) (m.Instruction, error)
	CreateInstruction(instruction m.Instruction) (m.Instruction, error)
	UpdateInstruction(instruction m.Instruction, recipeID uint) (m.Instruction, error)
	DeleteInstruction(recipeID uint) error
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
		h.utils.response500(w)
		return
	}

	data, err = h.recipeService.FindSingle(uint(rID))
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

func (h RecipeHandlers) GetIngredients(w http.ResponseWriter, r *http.Request, recipeID string) {
	var data []m.RecipeIngredient

	rID, err := strconv.Atoi(recipeID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	h.logger.Warnf("%s", "111")
	data, err = h.recipeService.FindRecipeIngredients(uint(rID))
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

func (h RecipeHandlers) Create(user *m.User, w http.ResponseWriter, r *http.Request) {
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

	recipe.Author.UserID = user.UserID

	data, err = h.recipeService.Create(recipe)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusCreated, data)
}

func (h RecipeHandlers) Update(w http.ResponseWriter, r *http.Request, recipeID string) {
	var recipe m.Recipe
	var data m.Recipe

	rID, err := strconv.Atoi(recipeID)
	if err != nil {
		h.utils.response500(w)
		return
	}

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

	data, err = h.recipeService.Update(recipe, uint(rID))
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

func (h RecipeHandlers) Delete(w http.ResponseWriter, r *http.Request, recipeID string) {
	rID, err := strconv.Atoi(recipeID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	if rID == 0 {
		h.utils.response400WithDetails(w, "ID is required")
		return
	}

	err = h.recipeService.Delete(uint(rID))
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.response204(w)
}

func (h RecipeHandlers) ImageUpload(w http.ResponseWriter, r *http.Request, recipeID string) {

	file, _, err := r.FormFile("file")
	if err != nil {
		h.utils.response400WithDetails(w, "bad request")
		return
	}

	rID, err := strconv.Atoi(recipeID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	_, err = h.recipeService.FindSingle(uint(rID))
	if err != nil {
		h.utils.response400WithDetails(w, "recipe does not exist")
		return
	}

	if s := h.imageService.UploadImage(file, rID); s {
		h.utils.response201(w)
		return
	}

	h.utils.response500(w)
}

func (h RecipeHandlers) GetInstruction(w http.ResponseWriter, r *http.Request, recipeID string) {
	var data m.Instruction

	rID, err := strconv.Atoi(recipeID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	data, err = h.recipeService.FindInstruction(uint(rID))
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

func (h RecipeHandlers) CreateInstruction(w http.ResponseWriter, r *http.Request) {
	var instruction m.Instruction
	var data m.Instruction

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	if err = json.Unmarshal(body, &instruction); err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	data, err = h.recipeService.CreateInstruction(instruction)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusCreated, data)
}

func (h RecipeHandlers) UpdateInstruction(w http.ResponseWriter, r *http.Request, recipeID string) {
	var instruction m.Instruction
	var data m.Instruction

	rID, err := strconv.Atoi(recipeID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	if err = json.Unmarshal(body, &instruction); err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	data, err = h.recipeService.UpdateInstruction(instruction, uint(rID))
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

func (h RecipeHandlers) DeleteInstruction(w http.ResponseWriter, r *http.Request, recipeID string) {
	rID, err := strconv.Atoi(recipeID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	err = h.recipeService.DeleteInstruction(uint(rID))
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.response204(w)
}
