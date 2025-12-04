package handlers

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/handlers/middleware"
	"github.com/stpotter16/biodata/internal/handlers/sessions"
	"github.com/stpotter16/biodata/internal/store"
)

func NewServer(store store.Store, sessionManager sessions.SessionManger) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, store)
	handler := middleware.PopulateSessionContext(sessionManager, mux)
	handler = middleware.LoggingWrapper(handler)
	return handler
}
