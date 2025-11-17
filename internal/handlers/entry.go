package handlers

import (
	"html/template"
	"net/http"
)

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
