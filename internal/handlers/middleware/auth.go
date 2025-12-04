package middleware

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/handlers/sessions"
)

func PopulateSessionContext(
	sessionManager sessions.SessionManger,
	next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCtx, err := sessionManager.PopulateSessionContext(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(sessionCtx))
	})
}

func NewViewAuthenticationRequiredMiddleware(sessionManager sessions.SessionManger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := sessionManager.SessionFromContext(r.Context())

			if err != nil {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
