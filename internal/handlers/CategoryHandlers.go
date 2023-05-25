package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	m "github.com/ihulsbus/cookbook/internal/models"
)

type CategoryService interface {
	FindAll() ([]m.Category, error)
	FindSingle(categoryID uint) (m.Category, error)
	Create(category m.Category) (m.Category, error)
	Update(category m.Category, categoryID uint) (m.Category, error)
	Delete(categoryID uint) error
}

type CategoryHandlers struct {
	categoryService CategoryService
	logger          LoggerInterface
	utils           HanderUtils
}

func NewCategoryHandlers(categorys CategoryService, logger LoggerInterface) *CategoryHandlers {
	return &CategoryHandlers{
		categoryService: categorys,
		logger:          logger,
		utils:           *NewHanderUtils(logger),
	}
}

func (h *CategoryHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	var data []m.Category

	data, err := h.categoryService.FindAll()
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

func (h *CategoryHandlers) Get(w http.ResponseWriter, r *http.Request, categoryID string) {
	var data m.Category

	tID, err := strconv.Atoi(categoryID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	if tID == 0 {
		h.utils.response400WithDetails(w, "ID is required")
		return
	}

	data, err = h.categoryService.FindSingle(uint(tID))
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

func (h *CategoryHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var data, category m.Category

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	if err = json.Unmarshal(body, &category); err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	data, err = h.categoryService.Create(category)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusCreated, data)
}

func (h *CategoryHandlers) Update(w http.ResponseWriter, r *http.Request, categoryID string) {
	var category, data m.Category

	tID, err := strconv.Atoi(categoryID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	if tID == 0 {
		h.utils.response400WithDetails(w, "ID is required")
		return
	}

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	if err = json.Unmarshal(body, &category); err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	data, err = h.categoryService.Update(category, uint(tID))
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

func (h *CategoryHandlers) Delete(w http.ResponseWriter, r *http.Request, categoryID string) {
	tID, err := strconv.Atoi(categoryID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	if tID == 0 {
		h.utils.response400WithDetails(w, "ID is required")
		return
	}

	err = h.categoryService.Delete(uint(tID))
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.response204(w)
}
