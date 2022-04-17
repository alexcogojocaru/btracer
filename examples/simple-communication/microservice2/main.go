package main

import (
	"context"
	"log"
	"net/http"

	"github.com/alexcogojocaru/btracer/exporters/bee"
	bhttp "github.com/alexcogojocaru/btracer/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

func ping(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ping route called"))
}

func main() {
	agentConfig := bee.AgentConfig{
		Host: "localhost",
		Port: 4576,
	}
	beeExporter, _ := bee.NewBeeExporter(&agentConfig)
	traceProvider := trace.NewTracerProvider(trace.WithBatcher(beeExporter))
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	http.Handle("/ping", bhttp.NewHandler(http.HandlerFunc(ping), "Ping"))
	http.ListenAndServe(":8090", nil)
}
