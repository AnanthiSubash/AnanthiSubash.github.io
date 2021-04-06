package router

import (
	"go-postgres/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/", middleware.IndexPage)
	router.HandleFunc("/login", middleware.LoginPage)
	router.HandleFunc("/logout", middleware.LogoutPage)
	router.HandleFunc("/admin", middleware.AdminPage)
	router.HandleFunc("/reviewer", middleware.ReviewerPage)
	router.HandleFunc("/dataentry", middleware.DataEntryPage)
	router.HandleFunc("/users", middleware.UsersPage)
	router.HandleFunc("/adduser", middleware.AddUser)
	router.HandleFunc("/edituser", middleware.EditUser)
	router.HandleFunc("/deleteuser", middleware.DeleteUser)
	router.HandleFunc("/showuser", middleware.ShowUser)
	router.HandleFunc("/questions", middleware.QuestionsPage)
	router.HandleFunc("/addquestion", middleware.AddQuestion)
	router.HandleFunc("/editquestion", middleware.EditQuestion)
	router.HandleFunc("/deletequestion", middleware.DeleteQuestion)
	router.HandleFunc("/showquestion", middleware.ShowQuestion)
	router.HandleFunc("/exams", middleware.ExamPage)
	router.HandleFunc("/addexam", middleware.AddExam)
	router.HandleFunc("/assignexam", middleware.AssignExam)
	router.HandleFunc("/dashboard", middleware.DashBoard)

	/*router.HandleFunc("/api/user/{id}", middleware.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user", middleware.GetAllUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newuser", middleware.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user/{id}", middleware.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteuser/{id}", middleware.DeleteUser).Methods("DELETE", "OPTIONS")
	*/
	return router
}
