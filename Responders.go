package Reminder

import (
	"encoding/json"
	"net/http"
)

func setupRequest(rWriter http.ResponseWriter) {
	rWriter.Header().Set("Server", "Go")
	rWriter.Header().Set("Access-Control-Allow-Origin", "*") //Leave this here for now (for development)
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
	encoder.Encode(note)
}

func return200Status(rWriter http.ResponseWriter) {
	setupRequest(rWriter)

	rWriter.WriteHeader(http.StatusOK)
}

func returnNotes(rWriter http.ResponseWriter, notes []*Note) {
	setupRequest(rWriter)

	encoder := json.NewEncoder(rWriter)
	encoder.Encode(notes)
}
