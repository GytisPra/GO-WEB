package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"web-app/internal/middleware"
	"web-app/internal/models"
	"web-app/internal/services"
	"web-app/pkg/utils"
)

type RegistrationHandler struct {
	userSerivce    *services.UserService
	accountService *services.AccountService
	sessionService *services.SessionService
}

func NewRegistrationHandler(sessionService *services.SessionService, userSerivce *services.UserService, accountService *services.AccountService) *RegistrationHandler {
	return &RegistrationHandler{sessionService: sessionService, userSerivce: userSerivce, accountService: accountService}
}

func (h *RegistrationHandler) ShowRegistrationForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, ok := middleware.FromContext(r.Context())
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err := r.Cookie("session_token")
	if err == nil {
		// user has a session cookie redirect him out of the registration
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	pageData := struct{}{}

	tmplData := utils.TemplateData{
		BuildTime:  utils.BuildTime,
		IsLoggedIn: false,
		ShowHeader: false,
		Page:       "register",
		Data:       pageData,
	}

	utils.RenderTemplate(w, "register.html", tmplData)
}

func (h *RegistrationHandler) Register(w http.ResponseWriter, r *http.Request) {
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
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.HandleError(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	var user *models.User
	user, err := h.userSerivce.GetUserByEmail(requestData.Email)
	if err != nil || user != nil {
		utils.HandleError(w, "An account with this email already exists", http.StatusBadRequest)
		return
	}

	user, err = h.userSerivce.CreateUser(requestData.Name, requestData.Email)
	if err != nil {
		utils.HandleError(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = h.accountService.CreateAccount(user.ID, "local", "local", nil, nil, nil, nil, nil, nil, nil, nil, nil, &requestData.Password)
	if err != nil {
		// Clean up user if account creation fails
		delErr := h.userSerivce.DeleteUserByID(user.ID)
		if delErr != nil {
			log.Printf("Failed to delete user after account creation error: %v", delErr)
		}
		utils.HandleError(w, fmt.Sprintf("Failed to create account: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
