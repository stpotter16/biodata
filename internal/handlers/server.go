package handlers

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/handlers/middleware"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux)
	handler := middleware.LoggingWrapper(mux)
	return handler
}
