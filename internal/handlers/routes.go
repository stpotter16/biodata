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
	apiAuthRequired := middleware.NewApiAuthenticationRequiredMiddleware(sessionManager, authorizer)
	mux.Handle("GET /api/entries", apiAuthRequired(entriesGet(store)))
	mux.Handle("GET /api/entries/{date}", apiAuthRequired(entryGet()))
	mux.Handle("POST /api/entry", apiAuthRequired(entryPost(store)))
	mux.Handle("PUT /api/entries/{date}", apiAuthRequired(entryPut(store)))
}
