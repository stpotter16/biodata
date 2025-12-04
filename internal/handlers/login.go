package handlers

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/handlers/authorization"
	"github.com/stpotter16/biodata/internal/handlers/sessions"
)

func loginPost(authorizer authorization.Authorizer, sessionManager sessions.SessionManger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !authorizer.Authorize("foo") {
			http.Error(w, "Invalid login attempt", http.StatusBadRequest)
			return
		}
		// TODO - handle error
		sessionManager.CreateSession(w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}
