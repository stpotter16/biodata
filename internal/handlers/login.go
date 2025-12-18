package handlers

import (
	"log"
	"net/http"

	"github.com/stpotter16/biodata/internal/handlers/authorization"
	"github.com/stpotter16/biodata/internal/handlers/sessions"
	"github.com/stpotter16/biodata/internal/parse"
)

func loginPost(authorizer authorization.Authorizer, sessionManager sessions.SessionManger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := parse.ParseLoginPost(r)
		if err != nil {
			log.Printf("Invalid login request %+v: %v", request, err)
			http.Error(w, "Invalid login request", http.StatusBadRequest)
			return
		}
		if !authorizer.Authorize(request) {
			log.Printf("Invalid login attempt: %+v", request)
			http.Error(w, "Invalid login attempt", http.StatusBadRequest)
			return
		}
		if err = sessionManager.CreateSession(w, r); err != nil {
			http.Error(w, "Failed to login", http.StatusInternalServerError)
		}
	}
}

func loginDelete(sessionManger sessions.SessionManger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := sessionManger.DeleteSession(w, r); err != nil {
			http.Error(w, "Failed to logout", http.StatusInternalServerError)
		}
	}
}
