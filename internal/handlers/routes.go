package handlers

import "net/http"

func addRoutes(mux *http.ServeMux) {
	// views
	mux.HandleFunc("GET /{$}", indexGet())
	mux.HandleFunc("GET /entry/new", newEntryGet())
	mux.HandleFunc("GET /entry/{date}/edit", editEntryGet())
	mux.Handle("GET /static/", http.StripPrefix("/static/", serveStaticFiles()))

	// API
	mux.HandleFunc("GET /api/entries", entriesGet())
	mux.HandleFunc("GET /api/entries/{date}", entryGet())
	mux.HandleFunc("POST /api/entry", entryPost())
	mux.HandleFunc("PUT /api/entries/{date}", entryPut())
}
