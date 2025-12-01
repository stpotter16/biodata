package handlers

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/stpotter16/biodata/internal/store"
	"github.com/stpotter16/biodata/internal/types"
)

//go:embed templates
var templateFS embed.FS

func indexGet(store store.Store) http.HandlerFunc {
	t := template.Must(
		template.New("base.html").
			ParseFS(
				templateFS,
				"templates/layouts/base.html", "templates/pages/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO - handle error
		entries, err := store.GetEntries()
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not read existing entries: %v", err), http.StatusInternalServerError)
			return
		}
		t.Execute(w, struct {
			Entries []types.Entry
		}{
			Entries: entries,
		})
	}
}

func newEntryGet() http.HandlerFunc {
	t := template.Must(
		template.New("base.html").
			ParseFS(
				templateFS,
				"templates/layouts/base.html", "templates/pages/new_entry.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, nil)
	}
}

func editEntryGet() http.HandlerFunc {
	t := template.Must(
		template.New("base.html").
			ParseFS(
				templateFS,
				"templates/layouts/base.html", "templates/pages/edit_entry.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, nil)
	}
}
