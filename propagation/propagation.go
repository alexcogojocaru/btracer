package propagation

import (
	"context"
	"log"
	"net/http"

	"github.com/alexcogojocaru/btracer/trace"
)

const TRACEPARENT_HEADER_TRACE = "Traceparent-TraceId"
const TRACEPARENT_HEADER_SPAN = "Traceparent-SpanId"
const TRACEPARENT_SEPARATOR = "-"

type Propagation interface {
	Inject(ctx context.Context, req *http.Request)
	Extract(ctx context.Context, header http.Header)
}

type Propagator struct{}

// Injects two headers containing the traceid and spanid of the current span
func (p *Propagator) Inject(ctx context.Context, req *http.Request) {
	contextHeader := ctx.Value("TraceHeader").(trace.ContextHeader)

	req.Header.Add(TRACEPARENT_HEADER_TRACE, contextHeader.ParentSpan.TraceID.ToString())
	req.Header.Add(TRACEPARENT_HEADER_SPAN, contextHeader.ParentSpan.SpanID.ToString())
}

func (p *Propagator) Extract(ctx context.Context, header http.Header) {
	traceId := header.Get(TRACEPARENT_HEADER_TRACE)
	spanId := header.Get(TRACEPARENT_HEADER_SPAN)

	log.Printf("%s %s", traceId, spanId)
}
