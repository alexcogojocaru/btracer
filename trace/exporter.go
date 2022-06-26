package trace

import (
	"context"
	"fmt"
	"log"

	agent "github.com/alexcogojocaru/btracer/proto-gen/btrace_proxy"
	"google.golang.org/grpc"
)

type Exporter struct {
	Client      agent.AgentClient
	ServiceName string
}

type AgentConfig struct {
	Host string
	Port int
}

func NewExporter(serviceName string, agentConfig AgentConfig) (*Exporter, error) {
	agentHost := fmt.Sprintf("%s:%d", agentConfig.Host, agentConfig.Port)

	conn, err := grpc.Dial(agentHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot dial %s", agentHost)
	}

	return &Exporter{
		Client:      agent.NewAgentClient(conn),
		ServiceName: serviceName,
	}, nil
}

func (e *Exporter) ExportSpan(ctx context.Context, span Span) {
	agentSpan := agent.Span{
		TraceID:      span.TraceID.ToString(),
		SpanID:       span.SpanID.ToString(),
		ParentSpanID: span.ParentSpanID.ToString(),
		Name:         span.Name,
		ServiceName:  e.ServiceName,
		Timestamp: &agent.Timestamp{
			Started:  span.StartTime.String(),
			Ended:    span.EndTime.String(),
			Duration: float32(span.Duration),
		},
	}

	log.Printf("%s parent=%s span=%s trace=%s\n", span.Name, span.ParentSpanID.ToString(), span.SpanID.ToString(), span.TraceID.ToString())

	e.Client.Send(ctx, &agentSpan)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}