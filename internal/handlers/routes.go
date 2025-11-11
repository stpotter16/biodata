package handlers

import "net/http"

func addRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", indexGet())
}
