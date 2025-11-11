package handlers

import (
	"fmt"
	"net/http"
)

func indexGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Printf("Hello from biodata server")
	}
}
