package main

import (
	"context"
	"log"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"

	bspan "github.com/alexcogojocaru/btracer/proto-gen/btrace_span"
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

func NormalizeSpan(span trace.ReadOnlySpan) bspan.Span {
	return bspan.Span{
		Name: span.Name(),
		CurrentContext: &bspan.Context{
			TraceID: span.SpanContext().TraceID().String(),
			SpanID:  span.SpanContext().SpanID().String(),
		},
		ParentContext: &bspan.Context{
			TraceID: span.Parent().TraceID().String(),
			SpanID:  span.Parent().SpanID().String(),
		},
	}
}

func NewExporter() (trace.SpanExporter, error) {
	return &Exporter{}, nil
}

func (e *Exporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	e.Mutex.Lock()
	for _, span := range spans {
		// log.Print(NormalizeSpan(span))
		log.Print(span.StartTime().Month().String())
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
