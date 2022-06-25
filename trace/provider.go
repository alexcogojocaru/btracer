package trace

import (
	"context"
)

const TRACE_HEADER = "TraceHeader"

type Provider interface {
	Start(ctx context.Context, name string) (context.Context, *Span)
	Shutdown() error
}

type TraceProvider struct {
	ServiceName     string
	Trace           Trace
	Encoder         Encoder
	KillSwitch      bool
	Channel         chan Span
	ShutdownChannel chan bool
	Exporter        *Exporter
}

type ContextHeader struct {
	Trace      Trace
	ParentSpan Span
	SpanName   string
}

type ContextSpan struct {
	TraceID  TID
	SpanID   SID
	SpanName string
}

func NewProvider(serviceName string) *TraceProvider {
	exporter, _ := NewExporter(
		serviceName,
		AgentConfig{
			Host: "localhost",
			Port: 4576,
		},
	)

	tp := &TraceProvider{
		ServiceName:     serviceName,
		Channel:         make(chan Span),
		ShutdownChannel: make(chan bool, 0),
		KillSwitch:      false,
		Exporter:        exporter,
	}

	go func() {
		tp.Stream()
	}()

	return tp
}

func (tp *TraceProvider) Start(ctx context.Context, name string) (context.Context, *Span) {
	span := &Span{
		Name: name,
	}

	if ctx.Value(TRACE_HEADER) == nil {
		traceToken := tp.Encoder.Compute(DEFAULT_TRACE_BYTES_SIZE)
		// copy(tp.Trace.TraceID[:], traceToken)

		// copy(span.ParentSpanID[:], NullSpanID[:])
		// copy(span.TraceID[:], tp.Trace.TraceID[:])

		copy(tp.Trace.TraceID[:], traceToken)
		copy(span.ParentSpanID[:], NullSpanID[:])
		copy(span.TraceID[:], traceToken)
	} else {
		// ctxValue := ctx.Value(TRACE_HEADER).(ContextHeader)
		// copy(span.ParentSpanID[:], ctxValue.ParentSpan.SpanID[:])
		// copy(span.TraceID[:], ctxValue.ParentSpan.TraceID[:])

		spanContext := ctx.Value(TRACE_HEADER).(ContextSpan)
		copy(span.ParentSpanID[:], spanContext.SpanID[:])
		copy(span.TraceID[:], spanContext.TraceID[:])
	}

	spanToken := tp.Encoder.Compute(DEFAULT_SPAN_BYTES_SIZE)
	copy(span.SpanID[:], spanToken)

	// InjectIntoContext(&ctx, ContextHeader{
	// 	Trace:      tp.Trace,
	// 	ParentSpan: *span,
	// 	SpanName:   name,
	// })

	InjectIntoContext(&ctx, ContextSpan{
		TraceID:  tp.Trace.TraceID,
		SpanID:   span.SpanID,
		SpanName: name,
	})

	tp.Channel <- *span
	return ctx, span
}

func (tp *TraceProvider) Stream() {
	for {
		span := <-tp.Channel
		tp.Exporter.ExportSpan(context.Background(), span)

		if tp.KillSwitch == true {
			tp.ShutdownChannel <- true
		}
	}
}

func (tp *TraceProvider) Shutdown() error {
	tp.KillSwitch = true
	<-tp.ShutdownChannel
	return nil
}

// func InjectIntoContext(ctx *context.Context, contextHeader ContextHeader) {
// 	*ctx = context.WithValue(*ctx, TRACE_HEADER, contextHeader)
// }

func InjectIntoContext(ctx *context.Context, contextHeader ContextSpan) {
	*ctx = context.WithValue(*ctx, TRACE_HEADER, contextHeader)
}
