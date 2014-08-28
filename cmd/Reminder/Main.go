package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Klathmon/Reminder-Go"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	//Init SQLite database
	dbConnection, err := sql.Open("sqlite3", "reminder.db")
	/*dbConnection.Exec(`
	CREATE TABLE Notes (
		ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		startDate DATETIME,
		dueDate DATETIME,
		nextDueDate DATETIME,
		done BOOLEAN,
		noteText TEXT);
	`)*/

	if err != nil {
		fmt.Println("Error opening DB connection: ", err)
		return
	}

	db, err := Reminder.PrepareStatements(dbConnection)
	if err != nil {
		fmt.Println("Error preparing statements: ", err)
		return
	}

	router := mux.NewRouter()

	router.HandleFunc("/Notes", func(rWriter http.ResponseWriter, req *http.Request) {
		Reminder.NewNote(rWriter, req, db)
	}).Methods("POST")

	router.HandleFunc("/Notes", func(rWriter http.ResponseWriter, req *http.Request) {
		Reminder.DeleteNote(rWriter, req, db)
	}).Methods("DELETE")

	router.HandleFunc("/Notes", func(rWriter http.ResponseWriter, req *http.Request) {
		Reminder.GetActiveNotes(rWriter, req, db)
	}).Methods("GET")

	http.Handle("/", router)

	http.ListenAndServe(":80", nil)
}
