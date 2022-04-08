package main

import (
	"context"
	"log"
	"net/http"

	"github.com/alexcogojocaru/btracer/exporters/bee"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

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
	// tracer := otel.Tracer("BTracer")

	// ctx, span := tracer.Start(context.Background(), "Main")
	// defer span.End()

	req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://localhost:8090/ping", nil)
	httpClient := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	_, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// resp, err := bhttp.Request(otelCtx, "http://localhost:8090/ping")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Print(resp)
}
