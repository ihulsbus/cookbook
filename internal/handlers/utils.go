package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h Handlers) NotImplemented(w http.ResponseWriter, r *http.Request) {
	h.response501(w)
}

func (h Handlers) response200(w http.ResponseWriter) {
	h.respondWithJSON(w, 200, "OK")
}

func (h Handlers) response201(w http.ResponseWriter) {
	h.respondWithJSON(w, 201, "Object created")
}

func (h Handlers) response204(w http.ResponseWriter) {
	h.respondWithJSON(w, 204, "Object deleted")
}

func (h Handlers) response501(w http.ResponseWriter) {
	h.respondWithError(w, 501, "Not Implemented")
}

func (h Handlers) response400WithDetails(w http.ResponseWriter, detail string) {
	h.respondWithError(w, 400, fmt.Sprintf("Bad Request. (%s)", detail))
}

func (h Handlers) response500WithDetails(w http.ResponseWriter, detail string) {
	h.respondWithError(w, 500, fmt.Sprintf("Internal Server Error. (%s)", detail))
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
