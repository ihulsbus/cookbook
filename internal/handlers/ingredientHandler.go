package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	c "github.com/ihulsbus/cookbook/internal/config"
	m "github.com/ihulsbus/cookbook/internal/models"
)

// Get all ingredients
func IngredientGetAll(w http.ResponseWriter, r *http.Request) {
	var data []m.Ingredient
	var responseCode int

	data, err := c.IngredientService.FindAllIngredients(c.IngredientService)
	if err != nil {
		response500WithDetails(w, err.Error())
	}

	responseCode = 200
	respondWithJSON(w, responseCode, data)
}

// Get a single ingredient
func IngredientGetSingle(w http.ResponseWriter, r *http.Request) {
	var data []m.Ingredient
	var responseCode int

	vars := mux.Vars(r)
	iID, err := strconv.Atoi(vars["ingredientID"])
	if err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	data, err = c.IngredientService.FindSingleIngredient(c.IngredientService, iID)
	if err != nil {
		response500WithDetails(w, err.Error())
	}

	responseCode = 200
	respondWithJSON(w, responseCode, data)
}

// Create an ingredient
func IngredientCreate(w http.ResponseWriter, r *http.Request) {
	var ingredient m.Ingredient
	var data m.Ingredient

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	body := buffer.String()

	if err = json.Unmarshal([]byte(body), &ingredient); err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	data, err = c.IngredientService.CreateIngredient(c.IngredientService, ingredient)
	if err != nil {
		response500WithDetails(w, err.Error())
	}

	respondWithJSON(w, 201, data)
}
