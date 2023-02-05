package handlers

import (
	"bytes"
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

// Get all ingredients
func (h Handlers) IngredientGetAll(w http.ResponseWriter, r *http.Request) {
	var data []m.Ingredient

	data, err := h.ingredientService.FindAll()
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, data)
}

// Get a single ingredient
func (h Handlers) IngredientGetSingle(w http.ResponseWriter, r *http.Request, ingredientID string) {
	var data m.Ingredient

	iID, err := strconv.Atoi(ingredientID)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	data, err = h.ingredientService.FindSingle(iID)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, data)
}

// Create an ingredient
func (h Handlers) IngredientCreate(w http.ResponseWriter, r *http.Request) {
	var ingredient m.Ingredient
	var data m.Ingredient

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	body := buffer.String()

	if err = json.Unmarshal([]byte(body), &ingredient); err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	data, err = h.ingredientService.Create(ingredient)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, data)
}

// TODO: Update handler

// Delete an ingredient
func (h Handlers) IngredientDelete(w http.ResponseWriter, r *http.Request) {
	var ingredient m.Ingredient

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	body := buffer.String()

	if err = json.Unmarshal([]byte(body), &ingredient); err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	if ingredient.ID == 0 {
		h.response400WithDetails(w, "ID is required")
		return
	}

	err = h.ingredientService.Delete(ingredient)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	h.response204(w)
}
