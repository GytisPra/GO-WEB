package auth

import (
	"fmt"
	"html/template"
	"net/http"
	"web-app/internal/services"
	"web-app/pkg/utils"
)

type LoginHandler struct {
	sessionService *services.SessionService
}

func NewLoginHandler(sessionService *services.SessionService) *LoginHandler {
	return &LoginHandler{sessionService: sessionService}
}

func (h *LoginHandler) LoginWithDiscord(w http.ResponseWriter, r *http.Request) {
	url := OAuth2Config.AuthCodeURL("state")

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *LoginHandler) loginWithLocal(w http.ResponseWriter, r *http.Request) {

}

func (h *LoginHandler) ShowLoginOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/login.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template file: %v", err), http.StatusInternalServerError)
		return
	}

	PageData := struct {
		BuildTime int64
	}{
		BuildTime: utils.BuildTime,
	}

	err = tmpl.ExecuteTemplate(w, "login.html", PageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}
