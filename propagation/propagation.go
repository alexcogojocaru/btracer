package propagation

import (
	"context"
	"net/http"

	"github.com/alexcogojocaru/btracer/trace"
)

const TRACEPARENT_HEADER_TRACE = "Traceparent-TraceId"
const TRACEPARENT_HEADER_SPAN = "Traceparent-SpanId"
const TRACEPARENT_HEADER_SERVICE = "Traceparent-Service"
const TRACEPARENT_SEPARATOR = "-"

type Propagation interface {
	Inject(ctx context.Context, req *http.Request)
	Extract(ctx context.Context, header http.Header) context.Context
}

type Propagator struct{}

// Injects two headers containing the traceid and spanid of the current span
func (p *Propagator) Inject(ctx context.Context, req *http.Request) {
	contextHeader := ctx.Value("TraceHeader").(trace.ContextSpan)

	req.Header.Add(TRACEPARENT_HEADER_TRACE, contextHeader.TraceID)
	req.Header.Add(TRACEPARENT_HEADER_SPAN, contextHeader.SpanID)
	req.Header.Add(TRACEPARENT_HEADER_SERVICE, contextHeader.ServiceName)
}

func (p *Propagator) Extract(ctx context.Context, header http.Header) context.Context {
	traceId := header.Get(TRACEPARENT_HEADER_TRACE)
	spanId := header.Get(TRACEPARENT_HEADER_SPAN)
	service := header.Get(TRACEPARENT_HEADER_SERVICE)

	extractedContext := context.Background()
	trace.InjectIntoContext(&extractedContext, trace.ContextSpan{
		TraceID:     traceId,
		SpanID:      spanId,
		ServiceName: service,
	})

	return extractedContext
}
