package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"web-app/internal/middleware"
	"web-app/internal/models"
	"web-app/internal/services"
)

type HomeHandler struct {
	sessionService *services.SessionService
}

func NewHomeHandler(sessionService *services.SessionService) *HomeHandler {
	return &HomeHandler{sessionService: sessionService}
}

func (h *HomeHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/home.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template file: %v", err), http.StatusInternalServerError)
		return
	}

	var isLoggedIn bool

	type ResponseData struct {
		IsLoggedIn bool
		User       *models.User
	}

	user, ok := middleware.FromContext(r.Context())
	if !ok {
		isLoggedIn = false
	} else {
		isLoggedIn = true
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", ResponseData{User: user, IsLoggedIn: isLoggedIn})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when executing template file: %v", err), http.StatusInternalServerError)
		return
	}

}
