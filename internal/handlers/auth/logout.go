package auth

import (
	"net/http"
	"time"
	"web-app/internal/middleware"
	"web-app/internal/services"
)

type LogoutHandler struct {
	sessionService *services.SessionService
}

func NewLogoutHandler(sessionService *services.SessionService) *LogoutHandler {
	return &LogoutHandler{sessionService: sessionService}
}

func (h *LogoutHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, ok := middleware.FromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	h.sessionService.LogUserOut(cookie.Value)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
