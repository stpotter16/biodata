package handlers

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/handlers/authorization"
	"github.com/stpotter16/biodata/internal/handlers/sessions"
	"github.com/stpotter16/biodata/internal/parse"
)

func loginPost(authorizer authorization.Authorizer, sessionManager sessions.SessionManger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := parse.ParseLoginPost(r)
		if err != nil {
			http.Error(w, "Invalid login request", http.StatusBadRequest)
			return
		}
		if !authorizer.Authorize(request) {
			http.Error(w, "Invalid login attempt", http.StatusBadRequest)
			return
		}
		// TODO - handle error
		sessionManager.CreateSession(w)
	}
}
