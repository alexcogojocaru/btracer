package http

import (
	"net/http"

	"github.com/alexcogojocaru/btracer/propagation"
)

var _ http.Handler = &Handler{}

type Handler struct {
	Handler    http.HandlerFunc
	Operation  string
	Propagator propagation.Propagation
}

func NewHandler(handler http.HandlerFunc, operation string) http.Handler {
	return &Handler{
		Handler:    handler,
		Operation:  operation,
		Propagator: &propagation.Propagator{},
	}
}

// this is a middleware: every request passes through this function
func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.Propagator.Extract(req.Context(), req.Header)
	h.Handler.ServeHTTP(w, req)
}
