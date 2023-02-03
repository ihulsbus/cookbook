package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	m "github.com/ihulsbus/cookbook/internal/models"
)

type RecipeService interface {
	FindAllRecipes() ([]m.Recipe, error)
	FindSingleRecipe(recipeID int) (m.Recipe, error)
	CreateRecipe(recipe m.Recipe) (m.Recipe, error)
	UpdateRecipe(recipe m.Recipe) (m.Recipe, error)
	DeleteRecipe(recipe m.Recipe) error
}

type ImageService interface {
	UploadImage(file multipart.File, recipeID int) bool
}

func (h Handlers) RecipeGetAll(w http.ResponseWriter, r *http.Request) {
	var data []m.Recipe

	data, err := h.recipeService.FindAllRecipes()
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, data)
}

func (h Handlers) RecipeGet(w http.ResponseWriter, r *http.Request, recipeID string) {
	var data m.Recipe

	rID, err := strconv.Atoi(recipeID)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	data, err = h.recipeService.FindSingleRecipe(rID)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, data)
}

func (h Handlers) RecipeCreate(w http.ResponseWriter, r *http.Request) {
	var recipe m.Recipe
	var data m.Recipe

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	body := buffer.String()

	if err = json.Unmarshal([]byte(body), &recipe); err != nil {
		h.response400WithDetails(w, err.Error())
		return
	}

	data, err = h.recipeService.CreateRecipe(recipe)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, data)
}

func (h Handlers) RecipeImageUpload(w http.ResponseWriter, r *http.Request, recipeID string) {

	file, _, err := r.FormFile("file")
	if err != nil {
		h.response400WithDetails(w, "bad request")
		return
	}

	ID, err := strconv.Atoi(recipeID)
	if err != nil {
		h.response500(w)
		return
	}

	_, err = h.recipeService.FindSingleRecipe(ID)
	if err != nil {
		h.response400WithDetails(w, "recipe does not exist")
	}

	if s := h.imageService.UploadImage(file, ID); s {
		h.response201(w)
		return
	}

	h.response500(w)

}

func (h Handlers) RecipeUpdate(w http.ResponseWriter, r *http.Request) {
	var recipe m.Recipe
	var data m.Recipe

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	body := buffer.String()

	if err = json.Unmarshal([]byte(body), &recipe); err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	if recipe.ID == 0 {
		h.response400WithDetails(w, "ID is required")
		return
	}

	data, err = h.recipeService.UpdateRecipe(recipe)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, data)
}

func (h Handlers) RecipeDelete(w http.ResponseWriter, r *http.Request) {
	var recipe m.Recipe

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	body := buffer.String()

	if err = json.Unmarshal([]byte(body), &recipe); err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	if recipe.ID == 0 {
		h.response400WithDetails(w, "ID is required")
		return
	}

	err = h.recipeService.DeleteRecipe(recipe)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	h.response204(w)
}
