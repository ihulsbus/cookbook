package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h Handlers) response201(w http.ResponseWriter) {
	h.respondWithJSON(w, http.StatusCreated, "Object created")
}

func (h Handlers) response204(w http.ResponseWriter) {
	h.respondWithJSON(w, http.StatusNoContent, "Object deleted")
}

func (h Handlers) response500(w http.ResponseWriter) {
	h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
}

func (h Handlers) response400WithDetails(w http.ResponseWriter, detail string) {
	h.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Bad Request. (%s)", detail))
}

func (h Handlers) response500WithDetails(w http.ResponseWriter, detail string) {
	h.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error. (%s)", detail))
}

func (h Handlers) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]interface{}{"code": code, "msg": message})
}

func (h Handlers) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)

	h.logger.Debugf("Sending response %s", response)

	if err != nil {
		h.logger.Warnf("Error occurred while returning error response to client: %s", err)
	}
}
