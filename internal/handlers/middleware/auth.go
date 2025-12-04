package middleware

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/handlers/sessions"
)

func PopulateSessionContext(
	sessionManager sessions.SessionManger,
	next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
