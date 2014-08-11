package Reminder

import (
	"encoding/json"
	"io"
	"time"
)

/*
Note holds a single note.
For one-time notes, DueDate and NoteText need to be set.If Done is set to true
then the note has been "completed". The program should never automatically mark
a note as completed, no matter how far past due.

Recurring Notes:
Set the DueDate for the next time the note will occur (if every week then set
the due date to one week from now), and NextDueDate for the following time the
note will occur (if every week thenset this to one week after the DueDate)
When the note is retrieved the system should check if the DueDate has passed AND
the note is Done. If those are true then create a new Note with a DueDate of the
current note's NextDueDate, and a StartDate of the current note's DueDate.

If a `StartDate` is set, then the note should not display to the user
until the `StartDate` has passed.
*/
type Note struct {
	ID          int64     `json:"ID"`          //Must be globally unique
	StartDate   time.Time `json:"startDate"`   //Used for recurring notes
	DueDate     time.Time `json:"dueDate"`     //Exact time the Note is Due
	NextDueDate time.Time `json:"nextDueDate"` //Used for recurring notes
	Done        bool      `json:"done"`        //true == completed
	Text        string    `json:"text"`        //The note contents/text
}

//CreateNoteFromReader creates a new note from an io.Reader which has a JSON string.
func CreateNoteFromReader(reader io.Reader) (*Note, error) {
	note := &Note{}
	if err := json.NewDecoder(reader).Decode(note); err != nil {
		return &Note{}, err
	}

	return note, nil
}
