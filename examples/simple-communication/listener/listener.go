package main

import (
	"net/http"

	bhttp "github.com/alexcogojocaru/btracer/http"
	"github.com/alexcogojocaru/btracer/trace"
)

func ping(w http.ResponseWriter, req *http.Request) {
	provider := req.Context().Value("provider").(*trace.TraceProvider)
	ctx, span := provider.Start(req.Context(), "PingRoute")
	span.End()

	_, span2 := provider.Start(ctx, "PingSubRoute")
	span2.End()

	w.Write([]byte("It works"))
}

func main() {
	http.Handle("/ping", bhttp.NewHandler(ping, "Ping"))
	http.ListenAndServe("0.0.0.0:8090", nil)
}
