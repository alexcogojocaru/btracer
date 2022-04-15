package main

import (
	"context"
	"log"
	"net/http"

	"github.com/alexcogojocaru/btracer/exporters/bee"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"

	btrace_propagation "github.com/alexcogojocaru/btracer/propagation"
)

type SpanConfig struct {
	TraceID    string
	SpanID     string
	TraceFlags string
}

func main() {
	agentConfig := bee.AgentConfig{Host: "localhost", Port: 4576}
	beeExporter, _ := bee.NewBeeExporter(&agentConfig)

	traceProvider := trace.NewTracerProvider(trace.WithBatcher(beeExporter))
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer := otel.Tracer("BTracer")

	ctx, span := tracer.Start(context.Background(), "Main")
	defer span.End()

	req, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:8090/ping", nil)

	httpClient := &http.Client{
		// Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// currentSpan -> req.Header["Traceparent"] -> new generated span
	propagator := btrace_propagation.NewPropagator()

	propagator.Inject(ctx, httpClient)
	_, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	spanConfig, _ := propagator.Extract(ctx, req.Header)
	log.Print(spanConfig)
}
