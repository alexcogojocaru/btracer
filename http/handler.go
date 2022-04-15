package http

import (
	"log"
	"net/http"

	btrace_propagation "github.com/alexcogojocaru/btracer/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	otel_trace "go.opentelemetry.io/otel/trace"
)

var _ http.Handler = &Handler{}

type Handler struct {
	Handler    http.Handler
	Operation  string
	Exporter   trace.SpanExporter
	Propagator btrace_propagation.Propagator
}

func NewHandler(handler http.Handler, operation string) http.Handler {
	h := &Handler{
		Handler:   handler,
		Operation: operation,
		// Exporter:  exporter,
		Propagator: btrace_propagation.NewPropagator(),
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

// this is a middleware: every request passes through this function
func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	spanConfig, _ := h.Propagator.Extract(req.Context(), req.Header)

	spanContext := otel_trace.NewSpanContext(otel_trace.SpanContextConfig{
		TraceID:    spanConfig.TraceID,
		SpanID:     spanConfig.SpanID,
		TraceFlags: 0x1,
	})

	ctx := otel_trace.ContextWithSpanContext(req.Context(), spanContext)
	req = req.WithContext(ctx)
	log.Print(spanConfig)

	h.Handler.ServeHTTP(w, req)
}
