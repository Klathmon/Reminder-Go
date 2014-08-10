package Reminder

import (
	"database/sql"
	"encoding/json"
	"os"
	"testing"

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
	_, err := prepareStatements(setupMockDatabase())
	if err != nil {
		t.Error("Expected good result, got error (", err, ")")
	}
}

func TestPrepareBadQueries(t *testing.T) {
	t.Parallel()

	db, _ := sql.Open("sqlite3", ":memory:")

	database, err := prepareStatements(db)
	if err == nil {
		jsonstuff, _ := json.Marshal(database)
		t.Error("Expected error, got good result (", jsonstuff, ")")
	}
}

func TestCreateNoteGood(t *testing.T) {
	t.Parallel()

	database, _ := prepareStatements(setupMockDatabase())

	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")
	note, _ := CreateNoteFromReader(fileHandle)

	database.CreateNote(note)

	if note.ID != 1 {
		jsonstuff, _ := json.Marshal(note)
		t.Error("Note does not have an ID after being created in the database (", string(jsonstuff), ")")
	}
}
