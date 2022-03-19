package main

import (
	"context"
	"log"

	"github.com/alexcogojocaru/btracer/exporters/bee"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"

	bhttp "github.com/alexcogojocaru/btracer/http"
)

func main() {
	agentConfig := bee.AgentConfig{Host: "localhost", Port: 4576}
	beeExporter, _ := bee.NewBeeExporter(&agentConfig)

	log.Print(beeExporter)

	traceProvider := trace.NewTracerProvider(trace.WithBatcher(beeExporter))
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	otel.SetTracerProvider(traceProvider)
	tracer := otel.Tracer("BTracer")

	otelCtx, span := tracer.Start(context.Background(), "Main")
	defer span.End()

	resp, err := bhttp.Request("http://localhost:8090/ping", otelCtx)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(resp)
}
