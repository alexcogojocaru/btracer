package trace

import (
	"context"
)

const TRACE_HEADER = "TraceHeader"

type Provider interface {
	Start(ctx context.Context, name string) (context.Context, *Span)
	Shutdown(ctx context.Context) error
}

type TraceProvider struct {
	Trace     Trace
	Encoder   Encoder
	SpanCount int64
}

type ContextHeader struct {
	Trace      Trace
	ParentSpan Span
	SpanName   string
}

func NewProvider() (TraceProvider, error) {
	provider := TraceProvider{}

	token := provider.Encoder.Compute(DEFAULT_TRACE_BYTES_SIZE)
	copy(provider.Trace.TraceID[:], token)

	return provider, nil
}

func (tp *TraceProvider) Start(ctx context.Context, name string) (context.Context, *Span) {
	span := &Span{Name: name}

	spanToken := tp.Encoder.Compute(DEFAULT_SPAN_BYTES_SIZE)
	copy(span.SpanID[:], spanToken)

	ctx = context.WithValue(ctx, TRACE_HEADER, &ContextHeader{
		Trace:      tp.Trace,
		ParentSpan: *span,
		SpanName:   name,
	})

	return ctx, span
}

func (tp *TraceProvider) Shutdown(ctx context.Context) error {
	return nil
}
