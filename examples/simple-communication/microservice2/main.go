package main

import (
	"net/http"

	bhttp "github.com/alexcogojocaru/btracer/http"
)

func ping(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("dasdasdad"))
}

func main() {
	http.Handle("/ping", bhttp.NewHandler(http.HandlerFunc(ping), "Ping"))
	http.ListenAndServe(":8090", nil)
}
