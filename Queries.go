package Reminder

import (
	"database/sql"
	"errors"
)

//Statements is a container struct for the Database.
//It holds the connection to the DB as well as all prepared satements used
//throughout the application.
type Statements struct {
	db           *sql.DB
	createNote   *sql.Stmt
	updateNote   *sql.Stmt
	retrieveNote *sql.Stmt
	deleteNote   *sql.Stmt
}

func prepareStatements(sqlDBHandle *sql.DB) (*Statements, error) {
	var database Statements
	var err error

	database.db = sqlDBHandle

	database.createNote, err = database.db.Prepare("INSERT INTO Notes (startDate, dueDate, nextDueDate, done, noteText) VALUES (?, ?, ?, ?, ?)")
	database.retrieveNote, err = database.db.Prepare("SELECT * FROM Notes WHERE ID=?")
	database.updateNote, err = database.db.Prepare("UPDATE Notes SET startDate=?, dueDate=?, nextDueDate=?, done=?, noteText=? WHERE ID=?")
	database.deleteNote, err = database.db.Prepare("DELETE FROM Notes WHERE ID=?")
	if err != nil {
		return &Statements{}, err
	}

	return &database, nil
}

//CreateNote adds a new note to the system.
//This function modifies the given Note to add the ID returned from the database
func (database *Statements) CreateNote(note *Note) error {
	tx, _ := database.db.Begin()
	defer tx.Rollback()

	execResult, _ := tx.Stmt(database.createNote).Exec(note.StartDate, note.DueDate, note.NextDueDate, note.Done, note.Text)

	note.ID, _ = execResult.LastInsertId()

	tx.Commit()
	return nil
}

//RetrieveNote gets a single note from the database based on ID
func (database *Statements) RetrieveNote(ID int64) (*Note, error) {
	row := database.retrieveNote.QueryRow(ID)

	note := &Note{}

	err := row.Scan(&note.ID, &note.StartDate, &note.DueDate, &note.NextDueDate, &note.Done, &note.Text)
	if err != nil {
		return &Note{}, err
	}
	return note, nil
}

//UpdateNote updates a note already in the database.
//note requires a good ID at minimum.
func (database *Statements) UpdateNote(note *Note) error {
	tx, _ := database.db.Begin()
	defer tx.Rollback()

	execResult, err := tx.Stmt(database.updateNote).Exec(note.StartDate, note.DueDate, note.NextDueDate, note.Done, note.Text, note.ID)
	if err != nil {
		return err
	}
	rowsAffected, _ := execResult.RowsAffected()
	if rowsAffected != 1 {
		return errors.New("Update Affected more/less than one row (" + string(rowsAffected) + " rows affected)")
	}

	tx.Commit()
	return nil
}

//DeleteNote completely removes a note from the database.
//the passed note should not be used after this function
func (database *Statements) DeleteNote(note *Note) error {
	tx, _ := database.db.Begin()
	defer tx.Rollback()

	execResult, err := tx.Stmt(database.deleteNote).Exec(note.ID)
	if err != nil {
		return err
	}
	rowsAffected, _ := execResult.RowsAffected()
	if rowsAffected != 1 {
		return errors.New("Delete Affected more/less than one row (" + string(rowsAffected) + " rows affected)")
	}

	tx.Commit()
	return nil
}
