package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/stpotter16/biodata/internal/handlers/sessions"
	"github.com/stpotter16/biodata/internal/store"
	"github.com/stpotter16/biodata/internal/types"
)

type viewProps struct {
	CsrfToken string
}

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
		props := viewProps{CsrfToken: ""}
		if err := t.Execute(w, struct {
			viewProps
		}{
			viewProps: props,
		}); err != nil {
			log.Printf("Could not create login page: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
		}
	}
}

func indexGet(sessionManager sessions.SessionManger, store store.Store) http.HandlerFunc {
	t := template.Must(
		template.New("base.html").
			Funcs(funcMap).
			ParseFS(
				templateFS,
				"templates/layouts/base.html", "templates/pages/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		props, err := extractViewProps(sessionManager, r)
		if err != nil {
			http.Error(w, "Could not validate session", http.StatusBadRequest)
			return
		}
		entries, err := store.GetEntries(r.Context())
		if err != nil {
			log.Printf("Could not read existing entries: %v", err)
			http.Error(w, "Could not load entries - try again later", http.StatusInternalServerError)
			return
		}
		if err = t.Execute(w, struct {
			viewProps
			Entries []types.Entry
		}{
			viewProps: props,
			Entries:   entries,
		}); err != nil {
			log.Printf("Could not create entry page: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
		}
	}
}

func newEntryGet(sessionManager sessions.SessionManger) http.HandlerFunc {
	t := template.Must(
		template.New("base.html").
			Funcs(funcMap).
			ParseFS(
				templateFS,
				"templates/layouts/base.html", "templates/pages/new_entry.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		props, err := extractViewProps(sessionManager, r)
		if err != nil {
			http.Error(w, "Could not validate session", http.StatusBadRequest)
			return
		}
		if err := t.Execute(w, struct {
			viewProps
		}{
			viewProps: props,
		}); err != nil {
			log.Printf("Could not create new entry page: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
		}
	}
}

func editEntryGet(sessionManager sessions.SessionManger, store store.Store) http.HandlerFunc {
	t := template.Must(
		template.New("base.html").
			Funcs(funcMap).
			ParseFS(
				templateFS,
				"templates/layouts/base.html", "templates/pages/edit_entry.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		props, err := extractViewProps(sessionManager, r)
		if err != nil {
			http.Error(w, "Could not validate session", http.StatusBadRequest)
			return
		}
		entryDateStr := r.PathValue("date")
		entryDate, err := time.Parse("2006-01-02", entryDateStr)
		if err != nil {
			log.Printf("Could not parse entry date: %v", err)
			http.Error(w, "Invalid entry date", http.StatusBadRequest)
			return
		}
		entry, err := store.GetEntry(r.Context(), entryDate)
		if err != nil {
			log.Printf("Could not read existing entries: %v", err)
			http.Error(w, "Could not load entry data", http.StatusInternalServerError)
			return
		}
		if err = t.Execute(w, struct {
			viewProps
			Entry types.Entry
		}{
			viewProps: props,
			Entry:     entry,
		}); err != nil {
			log.Printf("Could not create edit entry page: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
		}
	}
}

func extractViewProps(s sessions.SessionManger, r *http.Request) (viewProps, error) {
	session, err := s.SessionFromContext(r.Context())
	if err != nil {
		return viewProps{}, err
	}
	props := viewProps{
		CsrfToken: session.CsrfToken,
	}
	return props, nil
}
