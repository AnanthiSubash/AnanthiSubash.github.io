package middleware

import (
	"log"
	"net/http"
)

func DbInsertHandler(r *http.Request, insertStmt string, formValue string) {
	db := createConnection()
	// close the db connection
	defer db.Close()
	r.ParseForm()
	_, err := db.Exec(insertStmt, r.FormValue(formValue))
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
}

func DbDeleteHandler(r *http.Request, deleteStmt string, id int) {
	db := createConnection()
	defer db.Close()
	_, err := db.Exec(deleteStmt, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
}
