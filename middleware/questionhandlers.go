package middleware

import (
	//"database/sql"
	// package to encode and decode the json into struct and vice versa
	//	"fmt"
	//	"go-postgres/models" // models package where User schema is defined
	"html/template"
	"log"
	"net/http" // used to access the request and response object of the api
	"strconv"
	//	"os"       // used to read the environment variable
	//	"strconv"  // package used to covert string into int type
	//	"github.com/gorilla/mux" // used to get the params from the route
	//	"html/template"
	//	"github.com/joho/godotenv" // package used to read the .env file
	//	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type Questions struct {
	ID       int64
	Qtype    string
	Question string
	UserID   int64
}

func QuestionsPage(w http.ResponseWriter, r *http.Request) {
	if getUserId(r) == "" {
		log.Printf("Not Logged In..")
		http.Redirect(w, r, "/", 302)
	}
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	err := templates.ExecuteTemplate(w, "questionsPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func AddQuestion(w http.ResponseWriter, r *http.Request) {
	// create the postgres db connection
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	err := templates.ExecuteTemplate(w, "addQuestionPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" {
		db := createConnection()
		// close the db connection
		defer db.Close()
		r.ParseForm()
		insertStmt := `insert into "questions"("question", "q_type","user_id") values($1, $2,$3)`

		_, err = db.Exec(insertStmt, r.FormValue("question"), r.FormValue("qtype"), getUserId(r))
		if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
		}
	}
}

func EditQuestion(w http.ResponseWriter, r *http.Request) {
	// create the postgres db connection
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	err := templates.ExecuteTemplate(w, "editQuestionPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" {
		db := createConnection()

		// close the db connection
		defer db.Close()
		r.ParseForm()
		updateStmt := `update "questions" set "question"=$1, "q_type"=$2,"user_id"=$3 where "questions_id"=$4`
		_, err = db.Exec(updateStmt, r.FormValue("question"), r.FormValue("qtype"), getUserId(r), r.FormValue("id"))
		if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
		}
	}
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {

	// create the postgres db connection

	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	err := templates.ExecuteTemplate(w, "deleteQuestionPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db := createConnection()

	// close the db connection
	defer db.Close()
	r.ParseForm()
	deleteStmt := `delete from "questions" where questions_id=$1`
	qid, _ := strconv.Atoi(r.FormValue("id"))
	log.Printf("%v,%T", qid, qid)
	db.Exec(deleteStmt, qid)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

}
func ShowQuestion(w http.ResponseWriter, r *http.Request) {
	var q []Questions
	// create the postgres db connection

	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	//err := templates.ExecuteTemplate(w, "showUserPage", nil)
	/*if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}*/
	db := createConnection()

	// close the db connection
	defer db.Close()
	rows, err := db.Query(`SELECT "questions_id","question","q_type","user_id" FROM "questions"`)
	if err != nil {
		log.Println(err)
		http.Error(w, "there was an error", http.StatusInternalServerError)
		return
	}
	var id int64
	var question string
	var qtype string
	var userid int64
	for rows.Next() {
		err = rows.Scan(&id, &question, &qtype, &userid)
		if err != nil {
			log.Println(err)
			http.Error(w, "there was an error", http.StatusInternalServerError)
			return
		}
		q = append(q, Questions{ID: id, Question: question, Qtype: qtype, UserID: userid})
		//return
	}

	log.Print(q)
	e := templates.ExecuteTemplate(w, "showQuestionPage", q)
	if e != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
}
