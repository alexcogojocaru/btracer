package main

import (
	"math/rand"
	"net/http"

	bhttp "github.com/alexcogojocaru/btracer/http"
	"github.com/alexcogojocaru/btracer/trace"
)

const STUB_URL = "http://localhost:8091/stub"

func ping(w http.ResponseWriter, req *http.Request) {
	provider := req.Context().Value("provider").(*trace.TraceProvider)
	ctx, span := provider.Start(req.Context(), "PingRoute")
	span.AddLog("ERROR", "Wrong Context")
	span.End()

	pingCtx, span2 := provider.Start(ctx, "PingSubRoute")
	span2.End()

	prob := rand.Float32()
	if prob >= 0.5 {
		client := bhttp.NewTraceClient()
		client.Request(pingCtx, "GET", STUB_URL, nil)
	}

	w.Write([]byte("It works"))
}

func main() {
	http.Handle("/ping", bhttp.NewHandler(ping, "Ping", "Listener_Microservice"))
	http.ListenAndServe("0.0.0.0:8090", nil)
}
