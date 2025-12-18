package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/stpotter16/biodata/internal/handlers/sessions"
	"github.com/stpotter16/biodata/internal/parse"
	"github.com/stpotter16/biodata/internal/store"
	"github.com/stpotter16/biodata/internal/types"
)

func entriesGet(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var apiEntries []types.EntryAPI
		entries, err := store.GetEntries(r.Context())
		if err != nil {
			log.Printf("Error loading entries: %v", err)
			http.Error(w, "Error loading entries", http.StatusInternalServerError)
			return
		}
		for _, entry := range entries {
			apiEntry := types.ToEntryApi(entry)
			apiEntries = append(apiEntries, apiEntry)
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(apiEntries); err != nil {
			log.Printf("Could not serialize payload %+v: %v", apiEntries, err)
			http.Error(w, "Could not serialize entries", http.StatusInternalServerError)
		}
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
		entry, err := store.GetEntry(r.Context(), entryDate)
		if err != nil {
			log.Printf("Could not load entry for date %s: %v", entryDateStr, err)
			http.Error(w, "Error loading entry", http.StatusInternalServerError)
			return
		}
		apiEntry := types.ToEntryApi(entry)

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(apiEntry); err != nil {
			log.Printf("Could not serialize payload %+v: %v", apiEntry, err)
			http.Error(w, "Could not serialize entry", http.StatusInternalServerError)
		}
	}
}

func entryPost(sessionManager sessions.SessionManger, store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := sessionManager.SessionFromContext(r.Context())
		if err != nil {
			http.Error(w, "Could not validate session", http.StatusBadRequest)
			return
		}
		if !validateSessionCsrf(session, r) {
			http.Error(w, "Invalid session token", http.StatusBadRequest)
			return
		}
		newEntry, err := parse.ParseEntryPost(r)
		if err != nil {
			http.Error(w, "Invalid new entry request", http.StatusBadRequest)
			return
		}

		if err = store.InsertEntry(r.Context(), newEntry); err != nil {
			log.Printf("Could not added new entry %+v: %v", newEntry, err)
			http.Error(w, "Could not add new entry", http.StatusInternalServerError)
		}
	}
}

func entryPut(sessionManager sessions.SessionManger, store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := sessionManager.SessionFromContext(r.Context())
		if err != nil {
			http.Error(w, "Could not validate session", http.StatusBadRequest)
			return
		}
		if !validateSessionCsrf(session, r) {
			http.Error(w, "Invalid session token", http.StatusBadRequest)
			return
		}
		updatedEntry, err := parse.ParseEntryPut(r)
		if err != nil {
			http.Error(w, "Invalid entry update request", http.StatusBadRequest)
			return
		}

		if err = store.UpdateEntry(r.Context(), updatedEntry); err != nil {
			log.Printf("Could not update entry %+v: %v", updatedEntry, err)
			http.Error(w, "Could not update entry", http.StatusInternalServerError)
		}
	}
}
