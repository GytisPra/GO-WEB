package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"web-app/internal/middleware"
	"web-app/internal/services"
	"web-app/pkg/utils"

	appErrors "web-app/pkg/errors"
)

type LoginHandler struct {
	sessionService *services.SessionService
}

func NewLoginHandler(sessionService *services.SessionService) *LoginHandler {
	return &LoginHandler{sessionService: sessionService}
}

func (h *LoginHandler) LoginWithDiscord(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.FromContext(r.Context())
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	url := OAuth2Config.AuthCodeURL("state")

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *LoginHandler) LoginWithLocal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, ok := middleware.FromContext(r.Context())
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.HandleError(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	session, err := h.sessionService.LoginWithLocal(requestData.Email, requestData.Password)
	if err != nil {
		if errors.Is(err, appErrors.ErrInvalidCredentials) {
			utils.HandleError(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		utils.HandleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.SessionToken,
		Expires:  session.Expires,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *LoginHandler) ShowLoginOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, ok := middleware.FromContext(r.Context())
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	pageData := struct{}{}

	tmplData := utils.TemplateData{
		BuildTime:  utils.BuildTime,
		IsLoggedIn: false,
		ShowHeader: false,
		Page:       "login",
		Data:       pageData,
	}

	utils.RenderTemplate(w, "login.html", tmplData)
}
