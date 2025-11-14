package handlers

import "net/http"

func addRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", indexGet())
	mux.Handle("GET /static/", http.StripPrefix("/static/", serveStaticFiles()))
}
