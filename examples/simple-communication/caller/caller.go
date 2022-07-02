package main

import (
	"context"

	thttp "github.com/alexcogojocaru/btracer/http"
	"github.com/alexcogojocaru/btracer/trace"
)

const LISTENER_URL = "http://localhost:8090/ping"

func main() {
	client := thttp.NewTraceClient()
	provider := trace.NewProvider("Caller_v1")
	defer provider.Shutdown()

	ctx, span := provider.Start(context.Background(), "Caller_Main")
	defer span.End()

	client.Request(ctx, "GET", LISTENER_URL, nil)
}
