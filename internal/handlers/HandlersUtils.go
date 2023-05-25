package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type LoggerInterface interface {
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

type HanderUtils struct {
	logger LoggerInterface
}

func NewHanderUtils(logger LoggerInterface) *HanderUtils {
	return &HanderUtils{
		logger: logger,
	}
}

func (h HanderUtils) response201(w http.ResponseWriter) {
	h.respondWithJSON(w, http.StatusCreated, "Object created")
}

func (h HanderUtils) response204(w http.ResponseWriter) {
	h.respondWithJSON(w, http.StatusNoContent, "Object deleted")
}

func (h HanderUtils) response404(w http.ResponseWriter) {
	h.respondWithError(w, http.StatusNotFound, "not found")
}

func (h HanderUtils) response500(w http.ResponseWriter) {
	h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
}

func (h HanderUtils) response400WithDetails(w http.ResponseWriter, detail string) {
	h.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Bad Request. (%s)", detail))
}

func (h HanderUtils) response500WithDetails(w http.ResponseWriter, detail string) {
	h.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error. (%s)", detail))
}

func (h HanderUtils) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]interface{}{"code": code, "msg": message})
}

func (h HanderUtils) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)

	h.logger.Debugf("Sending response %s", response)

	if err != nil {
		h.logger.Warnf("Error occurred while returning error response to client: %s", err)
	}
}

func (h HanderUtils) getBody(r io.ReadCloser) ([]byte, error) {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r)
	if err != nil {
		return nil, errors.New("unable to get body from request")
	}

	body := buffer.Bytes()

	return body, nil
}
