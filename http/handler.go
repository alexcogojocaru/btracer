package http

import (
	"log"
	"net/http"

	"go.opentelemetry.io/otel/sdk/trace"
)

var _ http.Handler = &Handler{}

type Handler struct {
	Handler   http.Handler
	Operation string
	Exporter  trace.SpanExporter
}

func NewHandler(handler http.Handler, operation string) http.Handler {
	h := &Handler{
		Handler:   handler,
		Operation: operation,
		// Exporter:  exporter,
	}

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Print(h.Operation)
	h.Handler.ServeHTTP(w, req)
}
