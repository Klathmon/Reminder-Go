package Reminder

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewNoteGood(t *testing.T) {
	t.Parallel()
	db, _ := prepareStatements(setupMockDatabase())
	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")
	testServer := httptest.NewServer(http.HandlerFunc(func(rWriter http.ResponseWriter, req *http.Request) {
		NewNote(rWriter, req, db)
	}))
	defer testServer.Close()

	resp, _ := http.Post(testServer.URL, "application/json", fileHandle)

	note, _ := CreateNoteFromReader(resp.Body)
	resp.Body.Close()

	if note.ID != 1 {
		jsonstuff, _ := json.Marshal(note)
		t.Error("Expected returned note to have ID of 1, got this instead (", string(jsonstuff), ")")
	}
}

func TestNewNoteBadRequest(t *testing.T) {
	t.Parallel()
	db, _ := prepareStatements(setupMockDatabase())
	fileHandle, _ := os.Open("./test_files/badJSONNote.json")
	testServer := httptest.NewServer(http.HandlerFunc(func(rWriter http.ResponseWriter, req *http.Request) {
		NewNote(rWriter, req, db)
	}))
	defer testServer.Close()

	resp, _ := http.Post(testServer.URL, "application/json", fileHandle)

	if resp.StatusCode != http.StatusBadRequest {
		jsonstuff, _ := json.Marshal(resp)
		t.Error("Expected 400 error, got (", string(jsonstuff), ")")
	}
}

func dontTestGetAllNotesGood(t *testing.T) {
	t.Parallel()
	db, _ := prepareStatements(setupMockDatabase())
	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")
	testServer := httptest.NewServer(http.HandlerFunc(func(rWriter http.ResponseWriter, req *http.Request) {
		GetActiveNotes(rWriter, req, db)
	}))
	defer testServer.Close()

	resp, _ := http.Post(testServer.URL, "application/json", fileHandle)

	note, _ := CreateNoteFromReader(resp.Body)
	resp.Body.Close()

	if note.ID != 1 {
		jsonstuff, _ := json.Marshal(note)
		t.Error("Expected returned note to have ID of 1, got this instead (", string(jsonstuff), ")")
	}
}
