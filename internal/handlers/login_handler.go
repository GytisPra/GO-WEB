package handlers

import (
	"net/http"
	"web-app/internal/middleware"
	"web-app/internal/services"
)

type LoginHandler struct {
	userService *services.UserService
}

func NewLoginHandler(userService *services.UserService) *LoginHandler {
	return &LoginHandler{userService: userService}
}

func (h *LoginHandler) LoginWithDiscordHandler(w http.ResponseWriter, r *http.Request) {
	url := middleware.OAuth2Config.AuthCodeURL("state")

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
