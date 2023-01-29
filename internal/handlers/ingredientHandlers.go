package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	m "github.com/ihulsbus/cookbook/internal/models"
)

// Get all ingredients
func (h Handlers) IngredientGetAll(w http.ResponseWriter, r *http.Request) {
	var data []m.Ingredient
	var responseCode int

	data, err := h.ingredientService.FindAllIngredients()
	if err != nil {
		h.response500WithDetails(w, err.Error())
	}

	responseCode = 200
	h.respondWithJSON(w, responseCode, data)
}

// Get a single ingredient
func (h Handlers) IngredientGetSingle(w http.ResponseWriter, r *http.Request) {
	var data []m.Ingredient
	var responseCode int

	vars := mux.Vars(r)
	iID, err := strconv.Atoi(vars["ingredientID"])
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	data, err = h.ingredientService.FindSingleIngredient(iID)
	if err != nil {
		h.response500WithDetails(w, err.Error())
	}

	responseCode = 200
	h.respondWithJSON(w, responseCode, data)
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

	data, err = h.ingredientService.CreateIngredient(ingredient)
	if err != nil {
		h.response500WithDetails(w, err.Error())
	}

	h.respondWithJSON(w, 201, data)
}

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

	err = h.ingredientService.DeleteIngredient(ingredient)
	if err != nil {
		h.response500WithDetails(w, err.Error())
	}

	h.response204(w)
}
