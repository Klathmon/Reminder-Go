package Reminder

import (
	"encoding/json"
	"errors"
	"net/http"
)

const (
	//DefaultReturn is the default number of notes to return in a GET request
	DefaultReturn = 100
)

func panicRecovery(rWriter http.ResponseWriter) {
	if r := recover(); r != nil {
		serverError(rWriter, errors.New("Server Panic!"))
	}
}

/*
NewNote creates a note in the system.
The request should contain a JSON representation of a note.

This will return a 200 status along with the fully completed note in the response.

If an error occurs a 4xx or 5xx status will be returned along with more
information about the error.
*/
func NewNote(rWriter http.ResponseWriter, req *http.Request, db *Statements) {
	defer panicRecovery(rWriter)

	note, err := CreateNoteFromReader(req.Body)
	if err != nil {
		badRequest(rWriter, err)
	}

	db.CreateNote(note)

	noteCreated(rWriter, note)
}

/*
GetActiveNotes gets all active notes within the parameters provided.
The request can contain a JSON key-value pair of the following:
startNumber: the number of results to skip before returning new results (used for pagination)
	defaults to 0
numberToReturn: the amount of notes you want returned in one pass.
	defaults to DefaultReturn
*/
func GetActiveNotes(rWriter http.ResponseWriter, req *http.Request, db *Statements) {
	defer panicRecovery(rWriter)

	decoder := json.NewDecoder(req.Body)
	var startNumber, numberToReturn int
	var params map[string]int

	if req.ContentLength != 0 {
		err := decoder.Decode(&params)
		if err != nil {
			badRequest(rWriter, err)
			return
		}
	}

	startNumber, exists := params["startNumber"]
	if !exists {
		startNumber = 0
	}

	numberToReturn, exists = params["numberToReturn"]
	if !exists {
		numberToReturn = DefaultReturn
	}

	notes, err := db.GetActiveNotes(startNumber, numberToReturn)
	if err != nil {
		serverError(rWriter, err)
	}

	returnNotes(rWriter, notes)

}
