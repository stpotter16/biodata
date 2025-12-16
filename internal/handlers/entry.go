package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/stpotter16/biodata/internal/parse"
	"github.com/stpotter16/biodata/internal/store"
	"github.com/stpotter16/biodata/internal/types"
)

func entriesGet(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var apiEntries []types.EntryAPI
		entries, err := store.GetEntries()
		for _, entry := range entries {
			// TODO - Do I really need to return an error here?
			apiEntry, _ := types.ToEntryApi(entry)
			apiEntries = append(apiEntries, apiEntry)
		}
		if err != nil {
			log.Printf("Error loading entries: %v", err)
			http.Error(w, "Error loading entries", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Rype", "application/json")
		json.NewEncoder(w).Encode(apiEntries)
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

		// TODO - handle error
		if err = store.InsertEntry(newEntry); err != nil {
			http.Error(w, fmt.Sprintf("Could not add new entry: %v", err), http.StatusInternalServerError)
		}
	}
}

func entryPut(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updatedEntry, err := parse.ParseEntryPut(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid entry update request: %v", err), http.StatusBadRequest)
			return
		}

		// TODO - handle error
		store.UpdateEntry(updatedEntry)
	}
}
