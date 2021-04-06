package middleware

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
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
	TemplateHandler(w, "indexPage")
}

func LogoutPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		clearSession(w)
		http.Redirect(w, r, "/", 302)
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
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
		http.Redirect(w, r, "/", 302)
	}
	TemplateHandler(w, "adminPage")
}

func ReviewerPage(w http.ResponseWriter, r *http.Request) {
	if getUserId(r) == "" {
		http.Redirect(w, r, "/", 302)
	}
	TemplateHandler(w, "reviewerPage")
}
func DataEntryPage(w http.ResponseWriter, r *http.Request) {
	if getUserId(r) == "" {
		http.Redirect(w, r, "/", 302)
	}
	TemplateHandler(w, "dataentryPage")
}

func UsersPage(w http.ResponseWriter, r *http.Request) {
	if getUserId(r) == "" {
		http.Redirect(w, r, "/", 302)
	}
	TemplateHandler(w, "usersPage")

	if r.Method == "POST" {
		r.ParseForm()
		db := createConnection()
		var password string
		var role int64
		var page string
		// close the db connection
		defer db.Close()
		stmt := `SELECT "password","role_id" FROM "users" where "user_name"=$1`
		err := db.QueryRow(stmt, r.FormValue("name")).Scan(&password, &role)
		if err != nil {
			log.Println(err)
			http.Error(w, "there was an error", http.StatusInternalServerError)
			return
		}
		userstmt := `SELECT "page" FROM "user_role" where "role_id"=$1`
		db.QueryRow(userstmt, role).Scan(&page)
		http.Redirect(w, r, page, 302)
	}
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	// create the postgres db connection
	TemplateHandler(w, "addUserPage")

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
	TemplateHandler(w, "editUserPage")

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
	TemplateHandler(w, "deleteUserPage")
	r.ParseForm()
	uid, _ := strconv.Atoi(r.FormValue("id"))
	DbDeleteHandler(r, `delete from "users" where user_id=$1`, uid)
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
