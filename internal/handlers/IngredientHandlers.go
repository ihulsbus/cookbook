package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	m "github.com/ihulsbus/cookbook/internal/models"
)

type IngredientService interface {
	FindAll() ([]m.Ingredient, error)
	FindSingle(ingredientID int) (m.Ingredient, error)
	Create(ingredient m.Ingredient) (m.Ingredient, error)
	Update(ingredient m.Ingredient) (m.Ingredient, error)
	Delete(ingredient m.Ingredient) error
}

type IngredientHandlers struct {
	ingredientService IngredientService
	logger            LoggerInterface
	utils             HanderUtils
}

func NewIngredientHandlers(ingredients IngredientService, logger LoggerInterface) *IngredientHandlers {
	return &IngredientHandlers{
		ingredientService: ingredients,
		logger:            logger,
		utils:             *NewHanderUtils(logger),
	}
}

// Get all ingredients
func (h IngredientHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	var data []m.Ingredient

	data, err := h.ingredientService.FindAll()
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

// Get a single ingredient
func (h IngredientHandlers) GetSingle(w http.ResponseWriter, r *http.Request, ingredientID string) {
	var data m.Ingredient

	iID, err := strconv.Atoi(ingredientID)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	data, err = h.ingredientService.FindSingle(iID)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

// Create an ingredient
func (h IngredientHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var ingredient m.Ingredient
	var data m.Ingredient

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
	}

	if err = json.Unmarshal(body, &ingredient); err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	data, err = h.ingredientService.Create(ingredient)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusCreated, data)
}

func (h IngredientHandlers) Update(w http.ResponseWriter, r *http.Request) {
	var ingredient m.Ingredient
	var data m.Ingredient

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
	}

	if err = json.Unmarshal(body, &ingredient); err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	if ingredient.ID == 0 {
		h.utils.response400WithDetails(w, "ID is required")
		return
	}

	data, err = h.ingredientService.Update(ingredient)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

// Delete an ingredient
func (h IngredientHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	var ingredient m.Ingredient

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
	}

	if err = json.Unmarshal(body, &ingredient); err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	if ingredient.ID == 0 {
		h.utils.response400WithDetails(w, "ID is required")
		return
	}

	err = h.ingredientService.Delete(ingredient)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.response204(w)
}
