package main

import (
	"context"
	"log"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"

	bagent "github.com/alexcogojocaru/btracer/proto-gen/btrace_agent"
)

type Exporter struct {
	Mutex  sync.Mutex
	Client bagent.AgentServiceClient
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

func NormalizeSpan(span trace.ReadOnlySpan) bagent.Span {
	return bagent.Span{
		Name: span.Name(),
		CurrentContext: &bagent.Context{
			TraceID: span.SpanContext().TraceID().String(),
			SpanID:  span.SpanContext().SpanID().String(),
		},
		ParentContext: &bagent.Context{
			TraceID: span.Parent().TraceID().String(),
			SpanID:  span.Parent().SpanID().String(),
		},
	}
}

func NewExporter() (trace.SpanExporter, error) {
	conn, err := grpc.Dial("localhost:4576", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Cannot dial localhost:4576")
	}

	client := bagent.NewAgentServiceClient(conn)

	return &Exporter{Client: client}, nil
}

func (e *Exporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	e.Mutex.Lock()

	var batch []*bagent.Span
	for _, span := range spans {
		var bSpan bagent.Span = NormalizeSpan(span)
		// log.Print(span.StartTime().Month().String())
		e.Client.StreamSpan(ctx, &bSpan)
		batch = append(batch, &bSpan)
	}

	e.Client.StreamBatch(ctx, &bagent.BatchSpan{Spans: batch})

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
