package utils

import (
	"html/template"
	"net/http"
)

var BuildTime int64

type TemplateData struct {
	BuildTime  int64
	IsLoggedIn bool
	Page       string
	ShowHeader bool
	Data       any
}

func RenderTemplate(w http.ResponseWriter, tmplName string, tmplData TemplateData) {
	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/"+tmplName)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "layout.html", tmplData)
	if err != nil {
		http.Error(w, "Execution error", http.StatusInternalServerError)
	}
}
