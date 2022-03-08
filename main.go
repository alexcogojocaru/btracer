package main

import (
	"context"
	"log"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Exporter struct {
	Mutex sync.Mutex
}

type Context struct {
	TraceID string
	SpanID  string
}

type Span struct {
	Name           string
	CurrentContext Context
	ParentContext  Context
}

func NormalizeSpan(readSpan trace.ReadOnlySpan) Span {
	return Span{
		Name: readSpan.Name(),
		CurrentContext: Context{
			TraceID: readSpan.SpanContext().TraceID().String(),
			SpanID:  readSpan.SpanContext().SpanID().String(),
		},
		ParentContext: Context{
			TraceID: readSpan.Parent().TraceID().String(),
			SpanID:  readSpan.Parent().SpanID().String(),
		},
	}
}

func NewExporter() (trace.SpanExporter, error) {
	return &Exporter{}, nil
}

func (e *Exporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	e.Mutex.Lock()
	for _, span := range spans {
		log.Print(NormalizeSpan(span))
	}
	e.Mutex.Unlock()

	return nil
}

func (e *Exporter) Shutdown(ctx context.Context) error {
	log.Print("Shutdown received")
	return nil
}

func main() {
	exporter, err := NewExporter()
	if err != nil {
		log.Fatal("Error when creating the exporter")
	}

	traceProvider := trace.NewTracerProvider(trace.WithBatcher(exporter))
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	otel.SetTracerProvider(traceProvider)
	tracer := otel.Tracer("BTracer")

	otelCtx, span := tracer.Start(context.Background(), "Main")
	_, span1 := tracer.Start(otelCtx, "SubMain")

	span1.End()
	span.End()
}
