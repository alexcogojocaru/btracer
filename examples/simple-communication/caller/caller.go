package main

import (
	"context"
	"net/http"

	"github.com/alexcogojocaru/btracer/propagation"
	"github.com/alexcogojocaru/btracer/trace"
)

const LISTENER_URL = "http://localhost:8090/ping"

func main() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", LISTENER_URL, nil)
	propagator := propagation.Propagator{}

	provider := trace.NewProvider("Caller_v1.0")
	defer provider.Shutdown()

	// provider.Start(context.Background(), "Caller_Main")
	ctx, _ := provider.Start(context.Background(), "Caller_Main")
	propagator.Inject(ctx, req)

	client.Do(req)

	// span := ctx.Value("TraceHeader").(trace.ContextHeader).ParentSpan.SpanID
	// log.Print(span.ToString())
}
