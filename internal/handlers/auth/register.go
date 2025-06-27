package auth

import (
	"net/http"
	"web-app/internal/services"
)

type RegisterHandler struct {
	sessionService *services.SessionService
}

func NewRegisterHandler(sessionService *services.SessionService) *RegisterHandler {
	return &RegisterHandler{sessionService: sessionService}
}

func (h *RegisterHandler) ShowRegistrationForm(w http.ResponseWriter, r *http.Request) {

}

func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {

}
