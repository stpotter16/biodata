package handlers

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/store"
)

func addRoutes(mux *http.ServeMux, store store.Store) {
	// views
	mux.HandleFunc("GET /{$}", indexGet(store))
	mux.HandleFunc("GET /entry/new", newEntryGet())
	mux.HandleFunc("GET /entry/{date}/edit", editEntryGet())
	mux.Handle("GET /static/", http.StripPrefix("/static/", serveStaticFiles()))

	// API
	mux.HandleFunc("GET /api/entries", entriesGet())
	mux.HandleFunc("GET /api/entries/{date}", entryGet())
	mux.HandleFunc("POST /api/entry", entryPost())
	mux.HandleFunc("PUT /api/entries/{date}", entryPut())
}
