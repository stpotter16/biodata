package handlers

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/handlers/middleware"
	"github.com/stpotter16/biodata/internal/store"
)

func NewServer(store store.Store) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, store)
	handler := middleware.PopulateSessionContext(mux)
	handler = middleware.LoggingWrapper(handler)
	return handler
}
