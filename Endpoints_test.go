package Reminder

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupEndpointTests() *Statements {
	db, _ := prepareStatements(setupMockDatabase())
	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")
	note, _ := CreateNoteFromReader(fileHandle)
	db.CreateNote(note) //ID = 1
	db.CreateNote(note) //ID = 2
	db.CreateNote(note) //ID = 3
	db.CreateNote(note) //ID = 4
	db.CreateNote(note) //ID = 5
	db.CreateNote(note) //ID = 6
	db.CreateNote(note) //ID = 7
	db.CreateNote(note) //ID = 8
	db.CreateNote(note) //ID = 9
	db.CreateNote(note) //ID = 10

	return db
}

func TestNewNoteGood(t *testing.T) {
	t.Parallel()
	db := setupEndpointTests()
	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")
	testServer := httptest.NewServer(http.HandlerFunc(func(rWriter http.ResponseWriter, req *http.Request) {
		NewNote(rWriter, req, db)
	}))
	defer testServer.Close()

	resp, _ := http.Post(testServer.URL, "application/json", fileHandle)

	note, _ := CreateNoteFromReader(resp.Body)
	resp.Body.Close()

	if note.ID != 11 {
		jsonstuff, _ := json.Marshal(note)
		t.Error("Expected returned note to have ID of 11, got this instead (", string(jsonstuff), ")")
	}
}

func TestNewNoteBadRequest(t *testing.T) {
	t.Parallel()
	db := setupEndpointTests()
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

func TestGetAllNotesGood(t *testing.T) {
	t.Parallel()

	db := setupEndpointTests()

	//Setup test server
	testServer := httptest.NewServer(http.HandlerFunc(func(rWriter http.ResponseWriter, req *http.Request) {
		GetActiveNotes(rWriter, req, db)
	}))
	defer testServer.Close()

	var notes []*Note

	//Test empty post
	dataToPost, _ := os.Open("./test_files/emptyFile.json")
	resp, _ := http.Post(testServer.URL, "", dataToPost)
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&notes)
	if err != nil {
		t.Fatal("Decoding Error: ", err)
	}
	if len(notes) != 10 {
		jsonstuff, _ := json.Marshal(notes)
		t.Error("Expected 10 notes returned, got (", len(notes), "), ", string(jsonstuff))
	}
	if notes[0].ID != 1 {
		t.Error("Expected ID 1 to be first response, got (", notes[0].ID, ")")
	}
	resp.Body.Close()

	//Test small limit
	dataToPost, _ = os.Open("./test_files/goodGetAllNotesSmallLimit.json")
	resp, _ = http.Post(testServer.URL, "", dataToPost)
	decoder = json.NewDecoder(resp.Body)
	err = decoder.Decode(&notes)
	if err != nil {
		jsonstuff, _ := json.Marshal(resp)
		t.Fatal("Decoding Error: ", err, string(jsonstuff))
	}
	if len(notes) != 5 {
		jsonstuff, _ := json.Marshal(notes)
		t.Error("Expected 10 notes returned, got (", len(notes), "), ", string(jsonstuff))
	}
	if notes[0].ID != 1 {
		t.Error("Expected ID 1 to be first response, got (", notes[0].ID, ")")
	}
	resp.Body.Close()

	//Test offset
	dataToPost, _ = os.Open("./test_files/goodGetAllNotesSmallLimit.json")
	resp, _ = http.Post(testServer.URL, "", dataToPost)
	decoder = json.NewDecoder(resp.Body)
	err = decoder.Decode(&notes)
	if err != nil {
		jsonstuff, _ := json.Marshal(resp)
		t.Fatal("Decoding Error: ", err, string(jsonstuff))
	}
	if len(notes) != 5 {
		jsonstuff, _ := json.Marshal(notes)
		t.Error("Expected 10 notes returned, got (", len(notes), "), ", string(jsonstuff))
	}
	if notes[0].ID != 1 {
		t.Error("Expected ID 1 to be first response, got (", notes[0].ID, ")")
	}
	resp.Body.Close()
}

func TestGetAllNotesBad(t *testing.T) {
	t.Parallel()

	db := setupEndpointTests()

	//Setup test server
	testServer := httptest.NewServer(http.HandlerFunc(func(rWriter http.ResponseWriter, req *http.Request) {
		GetActiveNotes(rWriter, req, db)
	}))
	defer testServer.Close()

	var notes []*Note

	dataToPost, _ := os.Open("./test_files/badGetAllNotes.json")
	resp, _ := http.Post(testServer.URL, "", dataToPost)
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&notes)
	if err == nil {
		jsonstuff, _ := json.Marshal(notes)
		t.Fatal("Expecting error, got good result: ", string(jsonstuff))
	}
	resp.Body.Close()
}
