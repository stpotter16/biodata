package handlers

import (
	"net/http"
)

func indexGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		response := []byte("Hello from biodata server")
		w.Write(response)
	}
}
