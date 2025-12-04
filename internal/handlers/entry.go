package handlers

import (
	"fmt"
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
