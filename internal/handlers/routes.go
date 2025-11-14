package handlers

import "net/http"

func addRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", indexGet())
	mux.HandleFunc("GET /entry/new", newEntryGet())
	mux.Handle("GET /static/", http.StripPrefix("/static/", serveStaticFiles()))
}
