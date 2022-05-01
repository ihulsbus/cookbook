package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	response501(w)
}

func response200(w http.ResponseWriter) {
	respondWithJSON(w, 200, "OK")
}

func response201(w http.ResponseWriter) {
	respondWithJSON(w, 201, "Object created")
}

func response204(w http.ResponseWriter) {
	respondWithJSON(w, 204, "Object deleted")
}

func response501(w http.ResponseWriter) {
	respondWithError(w, 501, "Not Implemented")
}

func response400WithDetails(w http.ResponseWriter, detail string) {
	respondWithError(w, 400, fmt.Sprintf("Bad Request. (%s)", detail))
}

func response500WithDetails(w http.ResponseWriter, detail string) {
	respondWithError(w, 500, fmt.Sprintf("Internal Server Error. (%s)", detail))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]interface{}{"code": code, "msg": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)

	log.Debugf("Sending response %s", response)

	if err != nil {
		log.Warnf("Error occurred while returning error response to client: %s", err)
	}
}
