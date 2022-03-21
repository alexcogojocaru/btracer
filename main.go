package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"

	bee "github.com/alexcogojocaru/btracer/exporters/bee"
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
	tracer := otel.Tracer("BTracer")

	otelCtx, span := tracer.Start(context.Background(), "Main")
	span.End()

	_, span2 := tracer.Start(otelCtx, "Main2")
	defer span2.End()
}
