package bee

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	bagent "github.com/alexcogojocaru/btracer/proto-gen/btrace_agent"
)

type BeeExporter struct {
	Client bagent.AgentServiceClient
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

func NewBeeExporter() (trace.SpanExporter, error) {
	conn, err := grpc.Dial("localhost:4576", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Cannot dial localhost:4576")
	}

	client := bagent.NewAgentServiceClient(conn)

	return &BeeExporter{Client: client}, nil
}

func (e *BeeExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	log.Print("Export Spans")

	for _, span := range spans {
		var bSpan bagent.Span = NormalizeSpan(span)

		// Send a grpc request with metadata injected in the request - https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md

		ctx = metadata.AppendToOutgoingContext(ctx, "key1", "val1", "key2", "val2")
		e.Client.StreamSpan(ctx, &bSpan)
	}

	return nil
}

func (e *BeeExporter) Shutdown(ctx context.Context) error {
	log.Print("Shutdown received")
	return nil
}
