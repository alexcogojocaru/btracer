package main

import (
	"context"

	thttp "github.com/alexcogojocaru/btracer/http"
	"github.com/alexcogojocaru/btracer/trace"
)

const LISTENER_URL = "http://localhost:8090/ping"

func main() {
	client := thttp.NewTraceClient()
	provider := trace.NewProvider("CallerHello")
	defer provider.Shutdown()

	ctx, span := provider.Start(context.Background(), "Caller_Main")
	defer span.End()
	span.AddLog("INFO", "Starting a main block")

	client.Request(ctx, "GET", LISTENER_URL, nil)

	_, cacheSpan := provider.Start(ctx, "CacheCall")
	defer cacheSpan.End()
	cacheSpan.AddLog("INFO", "redis cache call after the request")

	_, dbSpan := provider.Start(ctx, "DbCall")
	defer dbSpan.End()
}
