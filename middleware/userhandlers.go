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
type Users struct {
	ID       int64
	Username string
	Password string
	Role     int64
}
type Dashboard struct {
	ID       int64
	Username string
	Count    int64
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	log.Print("IndexPage Inside")
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	err := templates.ExecuteTemplate(w, "indexPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func LogoutPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		clearSession(w)
		http.Redirect(w, r, "/", 302)
	}

}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("Inside Login Handler")
	if r.Method == "POST" {
		log.Println("Inside POST")
		r.ParseForm()
		db := createConnection()
		var formName string
		var formPassword string
		var userid int64
		var password string
		var role int64
		var page string
		// close the db connection
		defer db.Close()
		formName = r.FormValue("name")
		formPassword = r.FormValue("password")
		if formName != "" && formPassword != "" {
			stmt := `SELECT "user_id","password","role_id" FROM "users" where "user_name"=$1`
			db.QueryRow(stmt, r.FormValue("name")).Scan(&userid, &password, &role)
			if formPassword != password {
				// handle login failure
			}
			userstmt := `SELECT "page" FROM "user_role" where "role_id"=$1`
			db.QueryRow(userstmt, role).Scan(&page)
			// Handle Session Cookie
			setSession(formName, strconv.Itoa(int(userid)), w)
			http.Redirect(w, r, page, 302)
		}
	}
}
func AdminPage(w http.ResponseWriter, r *http.Request) {

	if getUserId(r) == "" {
		log.Printf("Not Logged In..")
		http.Redirect(w, r, "/", 302)
	}

	log.Print("AdminPage Inside")
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	err := templates.ExecuteTemplate(w, "adminPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func ReviewerPage(w http.ResponseWriter, r *http.Request) {

	if getUserId(r) == "" {
		log.Printf("Not Logged In..")
		http.Redirect(w, r, "/", 302)
	}
	log.Print("AdminPage Inside")
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	err := templates.ExecuteTemplate(w, "reviewerPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func DataEntryPage(w http.ResponseWriter, r *http.Request) {

	if getUserId(r) == "" {
		log.Printf("Not Logged In..")
		http.Redirect(w, r, "/", 302)
	}
	log.Print("AdminPage Inside")
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	err := templates.ExecuteTemplate(w, "dataentryPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func UsersPage(w http.ResponseWriter, r *http.Request) {
	if getUserId(r) == "" {
		log.Printf("Not Logged In..")
		http.Redirect(w, r, "/", 302)
	}
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	err := templates.ExecuteTemplate(w, "usersPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" {
		log.Println("here")
		r.ParseForm()
		db := createConnection()
		var password string
		var role int64
		var page string
		// close the db connection
		defer db.Close()
		stmt := `SELECT "password","role_id" FROM "users" where "user_name"=$1`
		log.Print(r.FormValue("name"))
		db.QueryRow(stmt, r.FormValue("name")).Scan(&password, &role)
		if err != nil {
			log.Println(err)
			http.Error(w, "there was an error", http.StatusInternalServerError)
			return
		}
		userstmt := `SELECT "page" FROM "user_role" where "role_id"=$1`
		db.QueryRow(userstmt, role).Scan(&page)
		log.Printf("Redirecting to %v", page)
		http.Redirect(w, r, page, 302)
	}
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	// create the postgres db connection
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	err := templates.ExecuteTemplate(w, "addUserPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		db := createConnection()
		// close the db connection
		defer db.Close()
		r.ParseForm()
		insertStmt := `insert into "users"("user_name", "password","role_id") values($1, $2,$3)`

		rol, _ := strconv.Atoi(r.FormValue("role"))

		_, err := db.Exec(insertStmt, r.FormValue("name"), r.FormValue("password"), rol)
		if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
		}
	}
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	// create the postgres db connection
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	err := templates.ExecuteTemplate(w, "editUserPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" {
		db := createConnection()

		// close the db connection
		defer db.Close()
		r.ParseForm()
		updateStmt := `update "users" set "user_name"=$1, "password"=$2, "role_id"=$3 where "user_id"=$4`
		rol, _ := strconv.Atoi(r.FormValue("role"))
		uid, _ := strconv.Atoi(r.FormValue("id"))

		_, err := db.Exec(updateStmt, r.FormValue("name"), r.FormValue("password"), rol, uid)
		if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
		}
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// create the postgres db connection

	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	err := templates.ExecuteTemplate(w, "deleteUserPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db := createConnection()

	// close the db connection
	defer db.Close()
	r.ParseForm()
	deleteStmt := `delete from "users" where user_id=$1`
	uid, _ := strconv.Atoi(r.FormValue("id"))
	log.Printf("%v,%T", uid, uid)
	db.Exec(deleteStmt, uid)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

}
func ShowUser(w http.ResponseWriter, r *http.Request) {
	var u []Users
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	db := createConnection()

	defer db.Close()
	rows, err := db.Query(`SELECT "user_id","user_name","password","role_id" FROM "users"`)
	if err != nil {
		log.Println(err)
		http.Error(w, "there was an error", http.StatusInternalServerError)
		return
	}
	var id int64
	var username string
	var password string
	var role int64
	for rows.Next() {
		err = rows.Scan(&id, &username, &password, &role)
		if err != nil {
			log.Println(err)
			http.Error(w, "there was an error", http.StatusInternalServerError)
			return
		}
		u = append(u, Users{ID: id, Username: username, Password: password, Role: role})
		//return
	}

	log.Print(u)
	e := templates.ExecuteTemplate(w, "showUserPage", u)
	if e != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
}
func DashBoard(w http.ResponseWriter, r *http.Request) {
	var u []Dashboard
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	db := createConnection()

	defer db.Close()
	rows, err := db.Query(`SELECT "user_id","user_name" FROM "users"`)
	if err != nil {
		log.Println(err)
		http.Error(w, "there was an error", http.StatusInternalServerError)
		return
	}
	var id int64
	var username string
	var count int64
	for rows.Next() {
		err = rows.Scan(&id, &username)
		stmt := `SELECT count("user_id")from questions where "user_id"=$1`
		db.QueryRow(stmt, id).Scan(&count)
		if err != nil {
			log.Println(err)
			http.Error(w, "there was an error", http.StatusInternalServerError)
			return
		}
		u = append(u, Dashboard{ID: id, Username: username, Count: count})
		//return
	}

	log.Print(u)
	e := templates.ExecuteTemplate(w, "dashBoard", u)
	if e != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
}
