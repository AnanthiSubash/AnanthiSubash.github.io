package middleware

import (
	"html/template"
	"net/http"
)

func TemplateHandler(w http.ResponseWriter, page string) {
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("assets/*"))
	err := templates.ExecuteTemplate(w, page, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
