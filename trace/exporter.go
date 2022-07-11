package trace

import (
	"context"
	"fmt"
	"log"

	proxy "github.com/alexcogojocaru/btracer/proto-gen/btrace_proxy"
	"google.golang.org/grpc"
)

type Exporter struct {
	Client      proxy.ExporterClient
	ServiceName string
	Bypass      bool
}

type ConnectionDetails struct {
	Host string
	Port int
}

type ExporterConfig struct {
	Bypass          bool
	AgentConfig     ConnectionDetails
	CollectorConfig ConnectionDetails
}

func NewExporter(serviceName string, config ExporterConfig) (*Exporter, error) {
	// agentHost := fmt.Sprintf("%s:%d", agentConfig.Host, agentConfig.Port)

	// conn, err := grpc.Dial(agentHost, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("Cannot dial %s", agentHost)
	// }

	var host string
	if config.Bypass == false {
		host = fmt.Sprintf("%s:%d", config.AgentConfig.Host, config.AgentConfig.Port)
	} else {
		host = fmt.Sprintf("%s:%d", config.CollectorConfig.Host, config.CollectorConfig.Port)
	}

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot dial %s", host)
	}

	return &Exporter{
		Client:      proxy.NewExporterClient(conn),
		ServiceName: serviceName,
		Bypass:      config.Bypass,
	}, nil
}

func (e *Exporter) ExportSpan(ctx context.Context, span Span) {
	agentSpan := proxy.Span{
		TraceID:      span.TraceID.ToString(),
		SpanID:       span.SpanID.ToString(),
		ParentSpanID: span.ParentSpanID.ToString(),
		Name:         span.Name,
		ServiceName:  span.ServiceName,
		TraceService: span.TraceService,
		Timestamp: &proxy.Timestamp{
			Started:  span.StartTime.String(),
			Ended:    span.EndTime.String(),
			Duration: float32(span.Duration),
		},
		Logs: span.Logs,
	}

	log.Printf("%s parent=%s span=%s trace=%s\n", span.Name, span.ParentSpanID.ToString(), span.SpanID.ToString(), span.TraceID.ToString())

	e.Client.Send(ctx, &agentSpan)
}
