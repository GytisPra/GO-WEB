package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"web-app/internal/models"
	"web-app/internal/services"
	"web-app/pkg/utils"
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

	var isLoggedIn bool
	var user *models.User

	type ResponseData struct {
		IsLoggedIn bool
		User       models.User
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			isLoggedIn = false
			user = &models.User{}
		} else {
			utils.HandleError(w, fmt.Sprintf("Failed to retrieve session cookie: %v", err), http.StatusUnauthorized)
			return
		}
	} else {
		// Cookie found, check if the user is logged in
		isLoggedIn, user, err = h.sessionService.IsUserLoggedIn(cookie.Value)
		if err != nil {
			utils.HandleError(w, fmt.Sprintf("Failed to check if user is logged in: %v", err), http.StatusUnauthorized)
			return
		}
	}

	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/home.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template file: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", ResponseData{User: *user, IsLoggedIn: isLoggedIn})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when executing template file: %v", err), http.StatusInternalServerError)
		return
	}
}
