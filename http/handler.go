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

func NewHandlerFunc(fp func(http.ResponseWriter, *http.Request), operation string) http.Handler {
	h := &Handler{
		Handler:   http.HandlerFunc(fp),
		Operation: operation,
	}

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// this is a middleware: every request passes through this function
	log.Print(req.Header)
	h.Handler.ServeHTTP(w, req)
}
