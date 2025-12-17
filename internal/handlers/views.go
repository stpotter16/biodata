package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/stpotter16/biodata/internal/store"
	"github.com/stpotter16/biodata/internal/types"
)

//go:embed templates
var templateFS embed.FS

var funcMap = template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format("2006-01-02")
	},
}

func loginGet() http.HandlerFunc {
	t := template.Must(
		template.New("base.html").
			Funcs(funcMap).
			ParseFS(
				templateFS,
				"templates/layouts/base.html", "templates/pages/login.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		if err := t.Execute(w, nil); err != nil {
			log.Printf("Could not create login page: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
		}
	}
}

func indexGet(store store.Store) http.HandlerFunc {
	t := template.Must(
		template.New("base.html").
			Funcs(funcMap).
			ParseFS(
				templateFS,
				"templates/layouts/base.html", "templates/pages/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := store.GetEntries()
		if err != nil {
			log.Printf("Could not read existing entries: %v", err)
			http.Error(w, "Could not load entries - try again later", http.StatusInternalServerError)
			return
		}
		if err = t.Execute(w, struct {
			Entries []types.Entry
		}{
			Entries: entries,
		}); err != nil {
			log.Printf("Could not create entry page: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
		}
	}
}

func newEntryGet() http.HandlerFunc {
	t := template.Must(
		template.New("base.html").
			Funcs(funcMap).
			ParseFS(
				templateFS,
				"templates/layouts/base.html", "templates/pages/new_entry.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		if err := t.Execute(w, nil); err != nil {
			log.Printf("Could not create new entry page: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
		}
	}
}

func editEntryGet(store store.Store) http.HandlerFunc {
	t := template.Must(
		template.New("base.html").
			Funcs(funcMap).
			ParseFS(
				templateFS,
				"templates/layouts/base.html", "templates/pages/edit_entry.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		entryDateStr := r.PathValue("date")
		entryDate, err := time.Parse("2006-01-02", entryDateStr)
		if err != nil {
			log.Printf("Could not parse entry date: %v", err)
			http.Error(w, "Invalid entry date", http.StatusBadRequest)
			return
		}
		entry, err := store.GetEntry(entryDate)
		if err != nil {
			log.Printf("Could not read existing entries: %v", err)
			http.Error(w, "Could not load entry data", http.StatusInternalServerError)
			return
		}
		if err = t.Execute(w, struct {
			Entry types.Entry
		}{
			Entry: entry,
		}); err != nil {
			log.Printf("Could not create edit entry page: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
		}
	}
}
