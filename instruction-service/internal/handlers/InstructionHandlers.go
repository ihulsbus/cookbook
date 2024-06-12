package handlers

import (
	"encoding/json"
	m "instruction-service/internal/models"
	"net/http"

	"github.com/google/uuid"
)

type InstructionService interface {
	FindInstruction(recipeID uuid.UUID) (m.Instruction, error)
	CreateInstruction(instruction m.Instruction, recipeID uuid.UUID) (m.Instruction, error)
	UpdateInstruction(instruction m.Instruction, recipeID uuid.UUID) (m.Instruction, error)
	DeleteInstruction(recipeID uuid.UUID) error
}

type InstructionHandlers struct {
	ingredientService InstructionService
	logger            LoggerInterface
	utils             HanderUtils
}

func NewInstructionHandlers(service InstructionService, logger LoggerInterface) *InstructionHandlers {
	return &InstructionHandlers{
		ingredientService: service,
		logger:            logger,
		utils:             *NewHanderUtils(logger),
	}
}

func (h InstructionHandlers) GetInstruction(w http.ResponseWriter, r *http.Request, recipeID uuid.UUID) {
	var data m.Instruction

	data, err := h.ingredientService.FindInstruction(recipeID)
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

func (h InstructionHandlers) CreateInstruction(w http.ResponseWriter, r *http.Request, recipeID uuid.UUID) {
	var instruction m.Instruction
	var data m.Instruction

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	if err = json.Unmarshal(body, &instruction); err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	data, err = h.ingredientService.CreateInstruction(instruction, recipeID)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusCreated, data)
}

func (h InstructionHandlers) UpdateInstruction(w http.ResponseWriter, r *http.Request, recipeID uuid.UUID) {
	var instruction m.Instruction
	var data m.Instruction

	body, err := h.utils.getBody(r.Body)
	if err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	if err = json.Unmarshal(body, &instruction); err != nil {
		h.utils.response400WithDetails(w, err.Error())
		return
	}

	data, err = h.ingredientService.UpdateInstruction(instruction, recipeID)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.respondWithJSON(w, http.StatusOK, data)
}

func (h InstructionHandlers) DeleteInstruction(w http.ResponseWriter, r *http.Request, recipeID uuid.UUID) {
	err := h.ingredientService.DeleteInstruction(recipeID)
	if err != nil {
		h.utils.response500WithDetails(w, err.Error())
		return
	}

	h.utils.response204(w)
}
