package http

import (
	"context"
	"log"
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
	provider := trace.NewProvider("Listener", trace.ExporterConfig{})

	return &Handler{
		Handler:    handler,
		Operation:  operation,
		Propagator: &propagation.Propagator{},
		Provider:   provider,
	}
}

// this is a middleware: every request passes through this function
func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := h.Propagator.Extract(req.Context(), req.Header)
	log.Printf("Received %s name", ctx.Value(trace.TRACE_HEADER).(trace.ContextSpan).ServiceName)
	spanCtx, span := h.Provider.Start(ctx, h.Operation)

	ctx = context.WithValue(spanCtx, "provider", h.Provider)

	h.Handler.ServeHTTP(w, req.WithContext(ctx))
	span.End()
}
