package main

import (
	"context"
	"log"

	"github.com/alexcogojocaru/btracer/trace"
)

func main() {
	provider := trace.NewProvider("Caller")
	defer provider.Shutdown()

	ctx, _ := provider.Start(context.Background(), "Caller_Main")

	span := ctx.Value("TraceHeader").(trace.ContextHeader).ParentSpan.SpanID
	log.Print(span.ToString())
}
