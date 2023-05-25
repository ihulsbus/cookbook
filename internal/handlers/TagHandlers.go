package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	m "github.com/ihulsbus/cookbook/internal/models"
)

type TagService interface {
	FindAll() ([]m.Tag, error)
	FindSingle(tagID uint) (m.Tag, error)
	Create(tag m.Tag) (m.Tag, error)
	Update(tag m.Tag, tagID uint) (m.Tag, error)
	Delete(tagID uint) error
}

type TagHandlers struct {
	tagService TagService
	logger     LoggerInterface
	utils      HanderUtils
}

func NewTagHandlers(tags TagService, logger LoggerInterface) *TagHandlers {
	return &TagHandlers{
		tagService: tags,
		logger:     logger,
		utils:      *NewHanderUtils(logger),
	}
}

func (h *TagHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	var data []m.Tag

	data, err := h.tagService.FindAll()
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

func (h *TagHandlers) Get(w http.ResponseWriter, r *http.Request, tagID string) {
	var data m.Tag

	tID, err := strconv.Atoi(tagID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	if tID == 0 {
		h.utils.response400WithDetails(w, "ID is required")
		return
	}

	data, err = h.tagService.FindSingle(uint(tID))
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

func (h *TagHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var data, tag m.Tag

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	if err = json.Unmarshal(body, &tag); err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	data, err = h.tagService.Create(tag)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusCreated, data)
}

func (h *TagHandlers) Update(w http.ResponseWriter, r *http.Request, tagID string) {
	var tag, data m.Tag

	tID, err := strconv.Atoi(tagID)
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

	if err = json.Unmarshal(body, &tag); err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	data, err = h.tagService.Update(tag, uint(tID))
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

func (h *TagHandlers) Delete(w http.ResponseWriter, r *http.Request, tagID string) {
	tID, err := strconv.Atoi(tagID)
	if err != nil {
		h.utils.response500(w)
		return
	}

	if tID == 0 {
		h.utils.response400WithDetails(w, "ID is required")
		return
	}

	err = h.tagService.Delete(uint(tID))
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.response204(w)
}
