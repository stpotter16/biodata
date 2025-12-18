package handlers

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/handlers/sessions"
)

const CSRF_HEADER = "X-BIODATA-CSRF-TOKEN"

func validateSessionCsrf(session sessions.Session, r *http.Request) bool {
	expectedCsrfToken := session.CsrfToken
	csrfHeader := r.Header.Get(CSRF_HEADER)
	return expectedCsrfToken == csrfHeader
}
