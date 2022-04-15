package propagation

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	otel_trace "go.opentelemetry.io/otel/trace"
)

const TRACEPARENT_HEADER = "Traceparent"
const TRACEPARENT_SEPARATOR = "-"
const EMPTY_STRING = ""

type Propagator interface {
	Inject(ctx context.Context) (http.Client, error)
	Extract(ctx context.Context, req http.Header) (SpanConfig, error)
}

type propagator struct {
}

type SpanConfig struct {
	TraceID    otel_trace.TraceID
	SpanID     otel_trace.SpanID
	TraceFlags otel_trace.TraceFlags
}

func NewPropagator() *propagator {
	return &propagator{}
}

func (p *propagator) Inject(ctx context.Context) (http.Client, error) {
	return http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}, nil
}

func (p *propagator) Extract(ctx context.Context, req http.Header) (SpanConfig, error) {
	traceParentHeader := req.Get(TRACEPARENT_HEADER)
	if traceParentHeader == EMPTY_STRING {
		return SpanConfig{}, errors.New("Traceparent header missing from request")
	}

	headerInfo := strings.Split(traceParentHeader, TRACEPARENT_SEPARATOR)
	trid, err := otel_trace.TraceIDFromHex(headerInfo[1])
	if err != nil {
		return SpanConfig{}, err
	}

	spid, err := otel_trace.SpanIDFromHex(headerInfo[2])
	if err != nil {
		return SpanConfig{}, err
	}

	spanConfig := SpanConfig{
		TraceID: trid,
		SpanID:  spid,
	}

	return spanConfig, nil
}
