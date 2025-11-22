package handlers

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/handlers/middleware"
	"github.com/stpotter16/biodata/internal/store/sqlite"
)

func NewServer(store sqlite.Store) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, store)
	handler := middleware.LoggingWrapper(mux)
	return handler
}
