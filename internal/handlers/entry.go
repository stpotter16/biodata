package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/stpotter16/biodata/internal/parse"
	"github.com/stpotter16/biodata/internal/store"
	"github.com/stpotter16/biodata/internal/types"
)

func entriesGet(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var apiEntries []types.EntryAPI
		entries, err := store.GetEntries()
		if err != nil {
			log.Printf("Error loading entries: %v", err)
			http.Error(w, "Error loading entries", http.StatusInternalServerError)
			return
		}
		for _, entry := range entries {
			// TODO - Do I really need to return an error here?
			apiEntry, _ := types.ToEntryApi(entry)
			apiEntries = append(apiEntries, apiEntry)
		}

		w.Header().Set("Content-Rype", "application/json")
		json.NewEncoder(w).Encode(apiEntries)
	}
}

func entryGet(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entryDateStr := r.PathValue("date")
		entryDate, err := time.Parse("2006-01-02", entryDateStr)
		if err != nil {
			log.Printf("Error parsing path date %s in route: %v", entryDateStr, err)
			http.Error(w, "Invalid date", http.StatusBadRequest)
			return
		}
		entry, err := store.GetEntry(entryDate)
		if err != nil {
			log.Printf("Could not load entry for date %s: %v", entryDateStr, err)
			http.Error(w, "Error loading entry", http.StatusInternalServerError)
			return
		}
		// TODO - Do I need this error?
		apiEntry, _ := types.ToEntryApi(entry)

		w.Header().Set("Content-Rype", "application/json")
		json.NewEncoder(w).Encode(apiEntry)
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
