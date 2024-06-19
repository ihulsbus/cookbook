package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	m "ingredient-service/internal/models"
)

type IngredientService interface {
	FindAll() ([]m.IngredientDTO, error)
	FindUnits() ([]m.UnitDTO, error)
	FindSingle(ingredientDTO m.IngredientDTO) (m.IngredientDTO, error)
	Create(ingredientDTO m.IngredientDTO) (m.IngredientDTO, error)
	Update(ingredientDTO m.IngredientDTO) (m.IngredientDTO, error)
	Delete(ingredientDTO m.IngredientDTO) error
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
	var data []m.IngredientDTO

	data, err := h.ingredientService.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			h.utils.response404(w)
			return
		default:
			h.utils.response500WithDetails(w, err.Error())
			return
		}
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

// Get all units
func (h IngredientHandlers) GetUnits(w http.ResponseWriter, r *http.Request) {
	var data []m.UnitDTO

	data, err := h.ingredientService.FindUnits()
	if err != nil {
		switch err.Error() {
		case "not found":
			h.utils.response404(w)
			return
		default:
			h.utils.response500WithDetails(w, err.Error())
			return
		}
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

// Get a single ingredient
func (h IngredientHandlers) GetSingle(w http.ResponseWriter, r *http.Request, ingredientID string) {
	var data m.IngredientDTO

	iID, err := strconv.Atoi(ingredientID)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	data, err = h.ingredientService.FindSingle(uint(iID))
	if err != nil {
		switch err.Error() {
		case "not found":
			h.utils.response404(w)
			return
		default:
			h.utils.response500WithDetails(w, err.Error())
			return
		}
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

func (h IngredientHandlers) Update(w http.ResponseWriter, r *http.Request, ingredientID string) {
	var ingredient m.Ingredient
	var data m.Ingredient

	iID, err := strconv.Atoi(ingredientID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	if iID == 0 {
		h.utils.response400WithDetails(w, "ID is required")
		return
	}

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
	}

	if err = json.Unmarshal(body, &ingredient); err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	data, err = h.ingredientService.Update(ingredient, uint(iID))
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

// Delete an ingredient
func (h IngredientHandlers) Delete(w http.ResponseWriter, r *http.Request, ingredientID string) {

	iID, err := strconv.Atoi(ingredientID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	if iID == 0 {
		h.utils.response400WithDetails(w, "ID is required")
		return
	}

	err = h.ingredientService.Delete(uint(iID))
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.response204(w)
}
