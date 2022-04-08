package main

import (
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func ping(w http.ResponseWriter, req *http.Request) {
	span := trace.SpanFromContext(req.Context())
	defer span.End()

	log.Print(span)

	// w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("ping route called"))
}

func main() {
	// http.Handle("/ping", bhttp.NewHandler(http.HandlerFunc(ping), "Ping"))
	// http.Handle("/ping", bhttp.NewHandlerFunc(ping, "Ping"))

	handler := http.HandlerFunc(ping)
	wrappedHandler := otelhttp.NewHandler(handler, "/ping")

	http.Handle("/ping", wrappedHandler)
	http.ListenAndServe(":8090", nil)
}
