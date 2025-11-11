package handlers

import "net/http"

func NewServer() http.Handler {
	mux := http.NewServeMux()
	return mux
}
