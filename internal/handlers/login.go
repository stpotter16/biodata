package handlers

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/handlers/sessions"
)

func loginPost(sessionManager sessions.SessionManger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO - handle error
		sessionManager.CreateSession(w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}
