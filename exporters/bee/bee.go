package bee

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"

	bagent "github.com/alexcogojocaru/btracer/proto-gen/btrace_agent"
)

type BeeExporter struct {
	Client bagent.AgentServiceClient
}

type AgentConfig struct {
	Host string
	Port int
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
		Timestamp: &bagent.Timestamp{
			Started:  span.StartTime().String(),
			Ended:    span.EndTime().String(),
			Duration: 1,
		},
	}
}

func NewBeeExporter(agentConfig *AgentConfig) (trace.SpanExporter, error) {
	agentHost := fmt.Sprintf("%s:%d", agentConfig.Host, agentConfig.Port)

	conn, err := grpc.Dial(agentHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot dial %s", agentHost)
	}

	client := bagent.NewAgentServiceClient(conn)

	return &BeeExporter{Client: client}, nil
}

func (e *BeeExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	log.Print("Export Spans")

	for _, span := range spans {
		var bSpan bagent.Span = NormalizeSpan(span)

		// Send a grpc request with metadata injected in the request
		// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md

		// ctx = metadata.AppendToOutgoingContext(ctx, "key1", "val1", "key2", "val2")
		e.Client.StreamSpan(ctx, &bSpan)
	}

	return nil
}

func (e *BeeExporter) Shutdown(ctx context.Context) error {
	log.Print("Shutdown received")
	return nil
}
