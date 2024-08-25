package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/home.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template file: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when executing template file: %v", err), http.StatusInternalServerError)
		return
	}
}
