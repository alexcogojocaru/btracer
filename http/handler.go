package http

import (
	"net/http"

	"github.com/alexcogojocaru/btracer/propagation"
	"github.com/alexcogojocaru/btracer/trace"
)

var _ http.Handler = &Handler{}

type Handler struct {
	Handler    http.HandlerFunc
	Operation  string
	Propagator propagation.Propagation
	Provider   *trace.TraceProvider
}

func NewHandler(handler http.HandlerFunc, operation string) http.Handler {
	return &Handler{
		Handler:    handler,
		Operation:  operation,
		Propagator: &propagation.Propagator{},
		Provider:   trace.NewProvider("Listener"),
	}
}

// this is a middleware: every request passes through this function
func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// traceId, spanId := h.Propagator.Extract(req.Context(), req.Header)
	// h.Provider.Start()
	h.Handler.ServeHTTP(w, req)
}
