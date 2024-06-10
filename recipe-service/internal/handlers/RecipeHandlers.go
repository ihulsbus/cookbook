package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	m "recipe-service/internal/models"
)

type RecipeService interface {
	FindAll() ([]m.Recipe, error)
	FindSingle(recipeID uint) (m.Recipe, error)
	Create(recipe m.Recipe) (m.Recipe, error)
	Update(recipe m.Recipe, recipeID uint) (m.Recipe, error)
	Delete(recipeID uint) error
}

type RecipeHandlers struct {
	recipeService RecipeService
	logger        LoggerInterface
	utils         HanderUtils
}

func NewRecipeHandlers(recipes RecipeService, logger LoggerInterface) *RecipeHandlers {
	return &RecipeHandlers{
		recipeService: recipes,
		logger:        logger,
		utils:         *NewHanderUtils(logger),
	}
}

func (h RecipeHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	var data []m.Recipe

	data, err := h.recipeService.FindAll()
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

func (h RecipeHandlers) Get(w http.ResponseWriter, r *http.Request, recipeID string) {
	var data m.Recipe

	rID, err := strconv.Atoi(recipeID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	data, err = h.recipeService.FindSingle(uint(rID))
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

	recipe.AuthorID = user.UserID

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
