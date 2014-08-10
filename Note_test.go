package Reminder

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestCreateNoteFromReaderGoodNote(t *testing.T) {
	t.Parallel()
	fileHandle, _ := os.Open("./test_files/goodJSONNote.json")

	note, err := CreateNoteFromReader(fileHandle)
	if err != nil {
		t.Error("Expected good result, got error (", err, ")")
		return
	}

	if note.Text != "This is a test good note!" {
		noteOutput, _ := json.Marshal(note)
		t.Error("Did not get expected note text (", string(noteOutput), ")")
		return
	}
}
func TestCreateNoteFromReaderGoodNoteNoDoneField(t *testing.T) {
	t.Parallel()
	fileHandle, _ := os.Open("./test_files/goodJSONNoteNoDoneField.json")

	note, err := CreateNoteFromReader(fileHandle)
	if err != nil {
		t.Error("Expected good result, got error (", err, ")")
		return
	}

	if note.Text != "This is a test good note!" {
		noteOutput, _ := json.Marshal(note)
		t.Error("Did not get expected note text (", string(noteOutput), ")")
		return
	}
}
func TestCreateNoteFromReaderGoodNoteRepeating(t *testing.T) {
	t.Parallel()
	fileHandle, _ := os.Open("./test_files/goodJSONNoteRepeating.json")

	note, err := CreateNoteFromReader(fileHandle)
	if err != nil {
		t.Error("Expected good result, got error (", err, ")")
		return
	}

	testTime, _ := time.Parse(time.RFC3339Nano, "2014-08-02T01:00:00Z")

	if note.NextDueDate.Equal(testTime) {
		noteOutput, _ := json.Marshal(note)
		t.Error("Did not get expected note text (", string(noteOutput), ")")
		return
	}
}
func TestCreateNoteFromReaderBadNote(t *testing.T) {
	t.Parallel()
	fileHandle, _ := os.Open("./test_files/badJSONNote.json")

	note, err := CreateNoteFromReader(fileHandle)

	if err == nil {
		noteOutput, _ := json.Marshal(note)
		t.Error("Expected failure but a valid note was returned (", string(noteOutput), ")")
		return
	}
}
