package middleware

import (
	"context"
	"net/http"
	"web-app/internal/models"
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

			isLoggedin, user, err := sessionService.IsUserLoggedIn(cookie.Value)
			if err != nil || !isLoggedin {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Add user to context for downstream handlers
			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SoftAuthMiddleware(sessionService *services.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err == nil {
				isLoggedIn, user, err := sessionService.IsUserLoggedIn(cookie.Value)
				if err == nil && isLoggedIn {
					// Put user in context
					ctx := context.WithValue(r.Context(), userContextKey, user)
					r = r.WithContext(ctx)
				}
			}
			// Proceed regardless of login
			next.ServeHTTP(w, r)
		})
	}
}

func FromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(userContextKey).(*models.User)
	return user, ok
}
