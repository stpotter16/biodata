package handlers

import (
	"log"
	"net/http"
)

func loginPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received: %s", r.Body)
	}
}
