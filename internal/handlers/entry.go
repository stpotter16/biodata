package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/stpotter16/biodata/internal/parse"
	"github.com/stpotter16/biodata/internal/store"
)

func entriesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func entryGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func entryPost(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newEntry, err := parse.ParseEntryPost(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid new entry request: %v", err), http.StatusBadRequest)
			return
		}

		// TODO - handler error
		store.InsertEntry(newEntry)
	}
}

func entryPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := struct {
			Weight string `json:"weight"`
			Waist  string `json:"waist"`
			BP     string `json:"bp"`
		}{}

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&body); err != nil {
			log.Printf("Invalid edit entry request: %v", err)
			http.Error(w, fmt.Sprintf("Invalid edit entry request: %v", err), http.StatusBadRequest)
			return
		}

		log.Printf("Received %v", body)
	}
}
