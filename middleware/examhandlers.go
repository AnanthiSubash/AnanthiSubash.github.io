package middleware

import (
	"html/template"
	"log"
	"net/http"
)

type QExam struct {
	ID       int64
	Examname string
	QID      int64
	Question string
}

func AddExam(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, "addExamPage")
	if r.Method == "POST" {
		DbInsertHandler(r, `insert into "exams"("exam_name") values($1)`, "exam")
	}
}

func ExamPage(w http.ResponseWriter, r *http.Request) {
	if getUserId(r) == "" {
		http.Redirect(w, r, "/", 302)
	}
	TemplateHandler(w, "examPage")
}

func ShowAssignExam(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, "assignPage")
}

func AssignExam(w http.ResponseWriter, r *http.Request) {
	var qexam []QExam

	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))

	db := createConnection()
	// Exam Details
	rows, err := db.Query(`SELECT "exam_id","exam_name" FROM "exams"`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var id int64
	var examname string

	for rows.Next() {
		err = rows.Scan(&id, &examname)
		if err != nil {
			log.Println(err)
			http.Error(w, "there was an error", http.StatusInternalServerError)
			return
		}
		qexam = append(qexam, QExam{ID: id, Examname: examname})
	}

	// Question Details
	questionrow, err := db.Query(`SELECT "questions_id","question" FROM "questions" WHERE "exam_id" IS NULL`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var questionid int64
	var question string

	for questionrow.Next() {
		err = questionrow.Scan(&questionid, &question)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		qexam = append(qexam, QExam{ID: questionid, Question: question})
	}
	log.Print(qexam)

	e := templates.ExecuteTemplate(w, "assignExam", qexam)
	if e != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	if r.Method == "POST" {
		db := createConnection()
		// close the db connection
		defer db.Close()
		r.ParseForm()

	}
}
