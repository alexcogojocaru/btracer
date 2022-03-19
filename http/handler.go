package http

import (
	"log"
	"net/http"
)

var _ http.Handler = &Handler{}

type Handler struct {
	Handler   http.Handler
	Operation string
}

func NewHandler(handler http.Handler, operation string) http.Handler {
	h := &Handler{
		Handler:   handler,
		Operation: operation,
	}

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Print(h.Operation)
	h.Handler.ServeHTTP(w, req)
}
