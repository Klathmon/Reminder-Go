package Reminder

import "database/sql"

//Statements is a container struct for the Database.
//It holds the connection to the DB as well as all prepared satements used
//throughout the application.
type Statements struct {
	db         *sql.DB
	createNote *sql.Stmt
}

func prepareStatements(sqlDBHandle *sql.DB) (*Statements, error) {
	var database Statements
	var err error

	database.db = sqlDBHandle

	database.createNote, err = database.db.Prepare("INSERT INTO Notes (startDate, dueDate, nextDueDate, done, noteText) values(?, ?, ?, ?, ?)")
	if err != nil {
		return &Statements{}, err
	}

	return &database, nil
}

//CreateNote adds a new note to the system.
//This function modifies the given Note to add the ID returned from the database
func (database *Statements) CreateNote(note *Note) error {

	tx, err := database.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	execResult, err := tx.Stmt(database.createNote).Exec(note.StartDate, note.DueDate, note.NextDueDate, note.Done, note.Text)
	if err != nil {
		return err
	}

	note.ID, err = execResult.LastInsertId()
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}
