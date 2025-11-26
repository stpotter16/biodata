package handlers

import (
	"embed"
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
		entries, _ := store.GetEntries()
		t.Execute(w, struct {
			Entries []types.EntryDTO
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
