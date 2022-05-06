package trace_test

import (
	"context"
	"testing"

	"github.com/alexcogojocaru/btracer/trace"
)

func TestProviderCreation(t *testing.T) {
	tp := trace.NewProvider()

	spanName := "TestProviderCreation"
	ctx, span := tp.Start(context.Background(), spanName)

	contextSpanName := ctx.Value(trace.TRACE_HEADER).(*trace.ContextHeader).SpanName
	if contextSpanName != spanName && span.Name != spanName {
		t.Error("Different internal span name")
	}
}
