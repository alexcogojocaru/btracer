package main

import (
	"log"
	"net/http"

	bhttp "github.com/alexcogojocaru/btracer/http"
)

func ping(w http.ResponseWriter, req *http.Request) {
	log.Print(req.Context().Value("TraceHeader"))
	w.Write([]byte("It works"))
}

func main() {
	http.Handle("/ping", bhttp.NewHandler(ping, "Ping"))
	http.ListenAndServe("0.0.0.0:8090", nil)
}
