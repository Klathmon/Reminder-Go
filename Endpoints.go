package Reminder

import "net/http"

/*
NewNote creates a note in the system.
The request should contain a JSON representation of a note.

This will return a 200 status along with the fully completed note in the response.

If an error occurs a 4xx or 5xx status will be returned along with more
information about the error.
*/
func NewNote(rWriter http.ResponseWriter, req *http.Request, db *Statements) {

	note, err := CreateNoteFromReader(req.Body)
	if err != nil {
		badRequest(rWriter, err)
	}

	err = db.CreateNote(note)
	if err != nil {
		serverError(rWriter, err)
	}

	noteCreated(rWriter, note)
}
