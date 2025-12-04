package handlers

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/handlers/authorization"
	"github.com/stpotter16/biodata/internal/handlers/middleware"
	"github.com/stpotter16/biodata/internal/handlers/sessions"
	"github.com/stpotter16/biodata/internal/store"
)

func addRoutes(
	mux *http.ServeMux,
	store store.Store,
	sessionManager sessions.SessionManger,
	authorizer authorization.Authorizer) {
	// static
	mux.Handle("GET /static/", http.StripPrefix("/static/", serveStaticFiles()))

	// views
	mux.HandleFunc("GET /login", loginGet())

	// views that need authentication
	viewAuthRequired := middleware.NewViewAuthenticationRequiredMiddleware(sessionManager)
	mux.Handle("GET /{$}", viewAuthRequired(indexGet(store)))
	mux.Handle("GET /entry/new", viewAuthRequired(newEntryGet()))
	mux.Handle("GET /entry/{date}/edit", viewAuthRequired(editEntryGet(store)))

	// Auth
	mux.HandleFunc("POST /login", loginPost(authorizer, sessionManager))

	// API
	mux.HandleFunc("GET /api/entries", entriesGet())
	mux.HandleFunc("GET /api/entries/{date}", entryGet())
	mux.HandleFunc("POST /api/entry", entryPost(store))
	mux.HandleFunc("PUT /api/entries/{date}", entryPut(store))
}
