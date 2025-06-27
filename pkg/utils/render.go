package utils

import (
	"context"
	"html/template"
	"net/http"
	"web-app/internal/middleware"
)

var BuildTime int64

type BaseTemplateData struct {
	BuildTime  int64
	IsLoggedIn bool
}

type TemplateData struct {
	BuildTime  int64
	IsLoggedIn bool
	Data       any
}

func BuildTemplateData(ctx context.Context, pageData any) TemplateData {
	_, ok := middleware.FromContext(ctx)
	return TemplateData{
		BuildTime:  BuildTime,
		IsLoggedIn: ok,
		Data:       pageData,
	}
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
