package middleware

import (
	"context"
	"net/http"
	"time"
	"web-app/internal/services"
)

type key int

const userContextKey key = 0

func AuthMiddleware(sessionService *services.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			session, err := sessionService.GetSessionByToken(cookie.Value)
			if err != nil || session.Expires.Before(time.Now()) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Add user to context for downstream handlers
			ctx := context.WithValue(r.Context(), userContextKey, session.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SoftAuthMiddleware(sessionService *services.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err == nil {
				session, err := sessionService.GetSessionByToken(cookie.Value)
				if err != nil || session.Expires.Before(time.Now()) {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}
				ctx := context.WithValue(r.Context(), userContextKey, session.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func FromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userContextKey).(string)
	return userID, ok
}
