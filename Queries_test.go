package Reminder

import (
	"database/sql"
	"encoding/json"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	createNotesTable = `
    CREATE TABLE Notes (
      ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
      startDate DATETIME,
      dueDate DATETIME,
      nextDueDate DATETIME,
      done BOOLEAN,
      noteText TEXT);
    `
)

func setupMockDatabase() *sql.DB {
	db, _ := sql.Open("sqlite3", ":memory:")

	db.Exec(createNotesTable)

	return db
}

func TestPrepareGoodQueries(t *testing.T) {
	t.Parallel()
	_, err := PrepareStatements(setupMockDatabase())
	if err != nil {
		t.Error("Expected good result, got error (", err, ")")
	}
}

func TestPrepareBadQueries(t *testing.T) {
	t.Parallel()

	db, _ := sql.Open("sqlite3", ":memory:")

	database, err := PrepareStatements(db)
	if err == nil {
		jsonstuff, _ := json.Marshal(database)
		t.Error("Expected error, got good result (", jsonstuff, ")")
	}
}

func TestGetActiveNotesGood(t *testing.T) {
	t.Parallel()

	database, _ := PrepareStatements(setupMockDatabase())

	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")
	note, _ := CreateNoteFromReader(fileHandle)
	database.CreateNote(note) //ID = 1
	database.CreateNote(note) //ID = 2
	database.CreateNote(note) //ID = 3
	note.Done = true
	database.CreateNote(note) //ID = 4
	note.Done = false
	note.StartDate = time.Now().AddDate(0, 0, 1)
	database.CreateNote(note) //ID = 5

	notesRetrieved, err := database.GetActiveNotes(0, 100)
	if err != nil {
		t.Error("Expected good result, got error (", err, ")")
	}

	if len(notesRetrieved) != 3 {
		jsonstuff, _ := json.Marshal(notesRetrieved)
		t.Error("Expected 3 notes, got ", len(notesRetrieved), " (", string(jsonstuff), ")")
	}

	for _, value := range notesRetrieved {
		if value.ID != 1 && value.ID != 2 && value.ID != 3 {
			jsonstuff, _ := json.Marshal(value)
			t.Error("Got unexpected note returned (", string(jsonstuff), ")")
		}
	}
}

func TestCreateNoteGood(t *testing.T) {
	t.Parallel()

	database, _ := PrepareStatements(setupMockDatabase())

	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")
	note, _ := CreateNoteFromReader(fileHandle)

	database.CreateNote(note)

	if note.ID != 1 {
		jsonstuff, _ := json.Marshal(note)
		t.Error("Note does not have an ID after being created in the database (", string(jsonstuff), ")")
	}
}

func TestRetrieveNoteGood(t *testing.T) {
	t.Parallel()

	database, _ := PrepareStatements(setupMockDatabase())

	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")
	note, _ := CreateNoteFromReader(fileHandle)
	database.CreateNote(note)

	noteRetrieved, err := database.RetrieveNote(1)
	if err != nil {
		t.Error("Expected good result, got error (", err, ")")
	}

	if noteRetrieved.Text != "This is a test good note!" {
		jsonstuff, _ := json.Marshal(noteRetrieved)
		t.Error("Note retrieved from database does not have correct information (", string(jsonstuff), ")")
	}
}

func TestRetrieveNoteBadID(t *testing.T) {
	t.Parallel()

	database, _ := PrepareStatements(setupMockDatabase())

	//Create a note
	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")
	note, _ := CreateNoteFromReader(fileHandle)
	database.CreateNote(note)

	//Retrieve that note
	noteRetrieved, err := database.RetrieveNote(6)
	if err == nil {
		jsonstuff, _ := json.Marshal(noteRetrieved)
		t.Error("Expected error, got good result (", string(jsonstuff), ")")
	}
}

func TestUpdateNoteGood(t *testing.T) {
	t.Parallel()

	database, _ := PrepareStatements(setupMockDatabase())

	//Create a note
	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")
	note, _ := CreateNoteFromReader(fileHandle)
	database.CreateNote(note)

	//Update the note
	note.Text = "Updated Text!"
	err := database.UpdateNote(note)
	if err != nil {
		t.Error("Expected good result, got error (", err, ")")
	}

	//Retrieve the updated note to test
	noteRetrieved, err := database.RetrieveNote(1)
	if err != nil {
		t.Error("Expected good result, got error (", err, ")")
	}

	if noteRetrieved.Text != "Updated Text!" {
		jsonstuff, _ := json.Marshal(noteRetrieved)
		t.Error("Note retrieved from database does not have correct information (", string(jsonstuff), ")")
	}
}

func TestUpdateNoteBadID(t *testing.T) {
	t.Parallel()

	database, _ := PrepareStatements(setupMockDatabase())

	//Create a note
	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")
	note, _ := CreateNoteFromReader(fileHandle)
	database.CreateNote(note)

	//Attempt to update the note
	note.Text = "Updated Text!"
	note.ID = 5
	err := database.UpdateNote(note)
	if err == nil {
		jsonstuff, _ := json.Marshal(note)
		t.Error("Expected error, got good result (", string(jsonstuff), ")")
	}
}

func TestDeleteNoteGood(t *testing.T) {
	t.Parallel()

	database, _ := PrepareStatements(setupMockDatabase())

	//Create a note
	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")
	note, _ := CreateNoteFromReader(fileHandle)
	database.CreateNote(note)

	//Delete the note
	err := database.DeleteNote(note)
	if err != nil {
		t.Error("Expected good result, got error (", err, ")")
	}

	//attempt to retrieve the deleted note to test
	noteRetrieved, err := database.RetrieveNote(1)
	if err == nil {
		jsonstuff, _ := json.Marshal(noteRetrieved)
		t.Error("Expected error, got good result (", string(jsonstuff), ")")
	}
}

func TestDeleteNoteBadID(t *testing.T) {
	t.Parallel()

	database, _ := PrepareStatements(setupMockDatabase())

	//Create a note
	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")
	note, _ := CreateNoteFromReader(fileHandle)
	database.CreateNote(note)

	//Attempt to delete the note
	note.ID = 5
	err := database.DeleteNote(note)
	if err == nil {
		jsonstuff, _ := json.Marshal(note)
		t.Error("Expected error, got good result (", string(jsonstuff), ")")
	}
}
