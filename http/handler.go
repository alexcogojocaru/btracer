package http

import (
	"context"
	"log"
	"net/http"

	"github.com/alexcogojocaru/btracer/exporters/bee"
	btrace_propagation "github.com/alexcogojocaru/btracer/propagation"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	otel_trace "go.opentelemetry.io/otel/trace"
)

var _ http.Handler = &Handler{}

type Handler struct {
	Handler    http.Handler
	Operation  string
	Exporter   trace.SpanExporter
	Propagator btrace_propagation.Propagator
	Tracer     otel_trace.Tracer
}

func InitTraceProvider(ctx context.Context) {
	agentConfig := bee.AgentConfig{Host: "localhost", Port: 4576}
	beeExporter, _ := bee.NewBeeExporter(&agentConfig)

	traceProvider := trace.NewTracerProvider(trace.WithBatcher(beeExporter))
	defer func() {
		if err := traceProvider.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	otel.SetTracerProvider(traceProvider)
}

func NewHandler(handler http.Handler, operation string) http.Handler {
	// InitTraceProvider(context.Background())

	return &Handler{
		Handler:   handler,
		Operation: operation,
		// Exporter:  exporter,
		Propagator: btrace_propagation.NewPropagator(),
		Tracer:     otel.Tracer(""),
	}
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

	_, span := h.Tracer.Start(ctx, "Name")
	h.Handler.ServeHTTP(w, req)

	span.End()
}
