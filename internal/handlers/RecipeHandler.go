package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	c "github.com/ihulsbus/cookbook/internal/config"
	m "github.com/ihulsbus/cookbook/internal/models"
)

func RecipeGetAll(w http.ResponseWriter, r *http.Request) {
	var data []m.RecipeDTO
	var responseCode int

	data, err := c.RecipeService.FindAllRecipes(c.RecipeService)
	if err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	responseCode = 200
	respondWithJSON(w, responseCode, data)
}

func RecipeGet(w http.ResponseWriter, r *http.Request) {
	var data m.RecipeDTO
	var responseCode int

	vars := mux.Vars(r)
	rID, err := strconv.Atoi(vars["recipeID"])
	if err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	data, err = c.RecipeService.FindSingleRecipe(c.RecipeService, uint(rID))
	if err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	responseCode = 200
	respondWithJSON(w, responseCode, data)
}

func RecipeCreate(w http.ResponseWriter, r *http.Request) {
	var recipe m.Recipe
	var data m.RecipeDTO

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	body := buffer.String()

	if err = json.Unmarshal([]byte(body), &recipe); err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	data, err = c.RecipeService.CreateRecipe(c.RecipeService, recipe)
	if err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	respondWithJSON(w, 201, data)
}

func RecipeImageUpload(w http.ResponseWriter, r *http.Request) {
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
			response500WithDetails(w, err.Error())
		}

		uploadedFile.File = files[i]

		uploadedFiles = append(uploadedFiles, uploadedFile)
	}

	if err = c.RecipeService.UploadRecipeImages(c.RecipeService, uploadedFiles); err != nil {
		response500WithDetails(w, err.Error())
	}

	response201(w)
}

func RecipeUpdate(w http.ResponseWriter, r *http.Request) {
	var recipe m.Recipe
	var data m.RecipeDTO

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	body := buffer.String()

	if err = json.Unmarshal([]byte(body), &recipe); err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	if recipe.ID == 0 {
		response400WithDetails(w, "ID is required")
		return
	}

	data, err = c.RecipeService.UpdateRecipe(c.RecipeService, recipe)
	if err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	respondWithJSON(w, 200, data)
}

func RecipeDelete(w http.ResponseWriter, r *http.Request) {
	var recipe m.Recipe

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	body := buffer.String()

	if err = json.Unmarshal([]byte(body), &recipe); err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	if recipe.ID == 0 {
		response400WithDetails(w, "ID is required")
		return
	}

	err = c.RecipeService.DeleteRecipe(c.RecipeService, recipe)
	if err != nil {
		response500WithDetails(w, err.Error())
		return
	}

	response204(w)
}
