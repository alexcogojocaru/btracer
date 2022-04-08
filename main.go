package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// propagation.TraceContext{}.Extract(context.Background(), )

	// bag := propagation.Baggage{}
	// textMap := propagation.TextMapPropagator{}
	// bag.Inject(context.Background(), textMap)
	tracer := otel.Tracer("BTracer")

	otelCtx, span := tracer.Start(context.Background(), "Main")
	defer span.End()

	_, span1 := tracer.Start(otelCtx, "SubMain")
	defer span1.End()
}
