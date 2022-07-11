package trace

import (
	"context"
	"encoding/hex"
	"log"
)

const TRACE_HEADER = "TraceHeader"
const NO_AGENT = false

type Provider interface {
	Start(ctx context.Context, name string) (context.Context, *Span)
	Shutdown() error
}

type TraceProvider struct {
	ServiceName     string
	InitServiceName string
	Trace           Trace
	Encoder         Encoder
	KillSwitch      bool
	Channel         chan Span
	ShutdownChannel chan bool
	Exporter        *Exporter
	BypassAgent     bool
}

type ContextHeader struct {
	Trace      Trace
	ParentSpan Span
	SpanName   string
}

type ContextSpan struct {
	TraceID     string
	SpanID      string
	SpanName    string
	ServiceName string
}

func NewProvider(serviceName string, config ExporterConfig) *TraceProvider {
	if config == (ExporterConfig{}) {
		config = ExporterConfig{
			AgentConfig: ConnectionDetails{
				Host: "localhost",
				Port: 4576,
			},
			Bypass: false,
		}
	}

	exporter, _ := NewExporter(
		serviceName,
		config,
	)

	tp := &TraceProvider{
		ServiceName:     serviceName,
		InitServiceName: serviceName,
		Channel:         make(chan Span),
		ShutdownChannel: make(chan bool, 0),
		KillSwitch:      false,
		Exporter:        exporter,
		BypassAgent:     false,
	}

	go func() {
		tp.Stream()
	}()

	return tp
}

func (tp *TraceProvider) Start(ctx context.Context, name string) (context.Context, *Span) {
	span := &Span{
		Name:        name,
		Channel:     tp.Channel,
		ServiceName: tp.InitServiceName,
	}

	if ctx.Value(TRACE_HEADER) == nil {
		traceToken := tp.Encoder.Compute(DEFAULT_TRACE_BYTES_SIZE)
		copy(tp.Trace.TraceID[:], traceToken)
		copy(span.ParentSpanID[:], NullSpanID[:])
		copy(span.TraceID[:], traceToken)
	} else {
		spanContext := ctx.Value(TRACE_HEADER).(ContextSpan)
		span_id, _ := hex.DecodeString(spanContext.SpanID)
		trace_id, _ := hex.DecodeString(spanContext.TraceID)

		log.Printf("%s %s", spanContext.ServiceName, tp.ServiceName)
		tp.ServiceName = spanContext.ServiceName

		copy(tp.Trace.TraceID[:], trace_id)
		copy(span.ParentSpanID[:], span_id)
		copy(span.TraceID[:], trace_id)
	}

	spanToken := tp.Encoder.Compute(DEFAULT_SPAN_BYTES_SIZE)
	copy(span.SpanID[:], spanToken)

	ctxTraceID := hex.EncodeToString(tp.Trace.TraceID[:])
	ctxSpanID := hex.EncodeToString(span.SpanID[:])
	InjectIntoContext(&ctx, ContextSpan{
		TraceID:     ctxTraceID,
		SpanID:      ctxSpanID,
		SpanName:    name,
		ServiceName: tp.ServiceName,
	})

	span.TraceService = tp.ServiceName
	span.Start()
	return ctx, span
}

func (tp *TraceProvider) Stream() {
	for {
		// non-blocking channel fetch
		select {
		case span := <-tp.Channel:
			tp.Exporter.ExportSpan(context.Background(), span)
		default:
			// no message received
			if tp.KillSwitch == true {
				tp.ShutdownChannel <- true
			}
		}
	}
}

func (tp *TraceProvider) Shutdown() error {
	tp.KillSwitch = true

	for {
		select {
		case <-tp.ShutdownChannel:
			return nil
		}
	}
}

func InjectIntoContext(ctx *context.Context, contextHeader ContextSpan) {
	*ctx = context.WithValue(*ctx, TRACE_HEADER, contextHeader)
}
