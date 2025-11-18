package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func entriesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func entryGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func entryPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := struct {
			Date   string `json:"date"`
			Weight string `json:"weight"`
			Waist  string `json:"waist"`
			BP     string `json:"bp"`
		}{}

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&body); err != nil {
			log.Printf("Invalid auth request: %v", err)
			http.Error(w, fmt.Sprintf("Invalid new entry request: %v", err), http.StatusBadRequest)
			return
		}

		log.Printf("Received %v", body)
	}
}

func entryPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
