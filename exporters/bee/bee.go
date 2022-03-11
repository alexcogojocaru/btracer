package bee

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/sdk/trace"
)

type BeeExporter struct {
}

func NewBeeExporter() (trace.SpanExporter, error) {
	return &BeeExporter{}, nil
}

func (e *BeeExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	log.Print("Export Spans")
	return nil
}

func (e *BeeExporter) Shutdown(ctx context.Context) error {
	log.Print("Shutdown received")
	return nil
}
