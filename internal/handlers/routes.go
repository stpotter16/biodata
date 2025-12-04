package handlers

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/handlers/authorization"
	"github.com/stpotter16/biodata/internal/handlers/sessions"
	"github.com/stpotter16/biodata/internal/store"
)

func addRoutes(
	mux *http.ServeMux,
	store store.Store,
	sessionManager sessions.SessionManger,
	authorizer authorization.Authorizer) {
	// views
	mux.HandleFunc("GET /login", loginGet())
	mux.HandleFunc("GET /{$}", indexGet(store))
	mux.HandleFunc("GET /entry/new", newEntryGet())
	mux.HandleFunc("GET /entry/{date}/edit", editEntryGet(store))
	mux.Handle("GET /static/", http.StripPrefix("/static/", serveStaticFiles()))

	// Auth
	mux.HandleFunc("POST /login", loginPost(authorizer, sessionManager))

	// API
	mux.HandleFunc("GET /api/entries", entriesGet())
	mux.HandleFunc("GET /api/entries/{date}", entryGet())
	mux.HandleFunc("POST /api/entry", entryPost(store))
	mux.HandleFunc("PUT /api/entries/{date}", entryPut(store))
}
