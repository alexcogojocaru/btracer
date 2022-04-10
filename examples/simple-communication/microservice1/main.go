package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/alexcogojocaru/btracer/exporters/bee"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	otel_trace "go.opentelemetry.io/otel/trace"
)

type SpanConfig struct {
	TraceID    string
	SpanID     string
	TraceFlags string
}

// [version]-[trace-id]-[parent-id]-[trace-flags]
func Extract(traceparentHeader string) SpanConfig {
	split := strings.Split(traceparentHeader, " ")

	spanConfig := SpanConfig{
		TraceID:    split[1],
		SpanID:     split[2],
		TraceFlags: split[3],
	}

	return spanConfig
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

	httpClient := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	_, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// currentSpan -> req.Header["Traceparent"] -> new generated span
	log.Print(req.Header["Traceparent"][0]) // this is the intermediary span between 2 microservices

	extractedData := Extract(req.Header["Traceparent"][0])

	traceID, _ := otel_trace.TraceIDFromHex(extractedData.TraceID)
	spanID, _ := otel_trace.SpanIDFromHex(extractedData.SpanID)

	// traceID, _ := otel_trace.TraceIDFromHex(req.Header["Traceparent"][0])

	spanContext := otel_trace.NewSpanContext(otel_trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: otel_trace.FlagsSampled,
	})

	log.Print(spanContext)

	// contex := otel_trace.ContextWithSpanContext(ctx, spanContext)
	// _, span_1 := tracer.Start(contex, "Hello")
	// defer span_1.End()

	// log.Print(spanContext)

	// resp, err := bhttp.Request(otelCtx, "http://localhost:8090/ping")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Print(resp)
}
