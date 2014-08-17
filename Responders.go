package Reminder

import (
	"encoding/json"
	"net/http"
)

func setupRequest(rWriter http.ResponseWriter) {
	rWriter.Header().Set("Server", "Go")
}

func badRequest(rWriter http.ResponseWriter, err error) {
	setupRequest(rWriter)
	jsonEncoder := json.NewEncoder(rWriter)

	rWriter.WriteHeader(http.StatusBadRequest)
	jsonEncoder.Encode(err)
}

func serverError(rWriter http.ResponseWriter, err error) {
	setupRequest(rWriter)
	jsonEncoder := json.NewEncoder(rWriter)

	rWriter.WriteHeader(http.StatusInternalServerError)
	jsonEncoder.Encode(err)
}

func noteCreated(rWriter http.ResponseWriter, note *Note) {
	setupRequest(rWriter)

	encoder := json.NewEncoder(rWriter)

	rWriter.WriteHeader(http.StatusOK)
	encoder.Encode(note)
}

func returnNotes(rWriter http.ResponseWriter, notes []*Note) {
	setupRequest(rWriter)

	encoder := json.NewEncoder(rWriter)

	rWriter.WriteHeader(http.StatusOK)
	encoder.Encode(notes)
}
