package main

import (
	"net/http"

	bhttp "github.com/alexcogojocaru/btracer/http"
)

func ping(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ping route called"))
}

func main() {
	http.Handle("/ping", bhttp.NewHandler(http.HandlerFunc(ping), "Ping"))
	// http.Handle("/ping", bhttp.NewHandlerFunc(ping, "Ping"))

	// handler := http.HandlerFunc(ping)
	// wrappedHandler := otelhttp.NewHandler(handler, "/ping")

	// http.Handle("/ping", wrappedHandler)
	http.ListenAndServe(":8090", nil)
}
