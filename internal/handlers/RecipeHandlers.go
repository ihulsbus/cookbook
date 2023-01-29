package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	m "github.com/ihulsbus/cookbook/internal/models"
)

func (h Handlers) RecipeGetAll(w http.ResponseWriter, r *http.Request) {
	var data []m.RecipeDTO
	var responseCode int

	data, err := h.recipeService.FindAllRecipes()
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	responseCode = 200
	h.respondWithJSON(w, responseCode, data)
}

func (h Handlers) RecipeGet(w http.ResponseWriter, r *http.Request) {
	var data m.RecipeDTO
	var responseCode int

	vars := mux.Vars(r)
	rID, err := strconv.Atoi(vars["recipeID"])
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	data, err = h.recipeService.FindSingleRecipe(uint(rID))
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	responseCode = 200
	h.respondWithJSON(w, responseCode, data)
}

func (h Handlers) RecipeCreate(w http.ResponseWriter, r *http.Request) {
	var recipe m.Recipe
	var data m.RecipeDTO

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

	data, err = h.recipeService.CreateRecipe(recipe)
	if err != nil {
		h.response500WithDetails(w, err.Error())
		return
	}

	h.respondWithJSON(w, 201, data)
}

func (h Handlers) RecipeImageUpload(w http.ResponseWriter, r *http.Request) {
	var uploadedFiles []m.RecipeFile
	var err error
	vars := mux.Vars(r)

	err = r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	formdata := r.MultipartForm // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File["multiplefiles"] // grab the filenames

	for i := range files { // loop through the files one by one
		var uploadedFile m.RecipeFile

		uploadedFile.ID, err = strconv.Atoi(vars["recipeID"])
		if err != nil {
			h.response500WithDetails(w, err.Error())
		}

		uploadedFile.File = files[i]

		uploadedFiles = append(uploadedFiles, uploadedFile)
	}

	if err = h.recipeService.UploadRecipeImages(uploadedFiles); err != nil {
		h.response500WithDetails(w, err.Error())
	}

	h.response201(w)
}

func (h Handlers) RecipeUpdate(w http.ResponseWriter, r *http.Request) {
	var recipe m.Recipe
	var data m.RecipeDTO

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

	h.respondWithJSON(w, 200, data)
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
