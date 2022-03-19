package main

import (
	"net/http"
)

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("custom"))
}

func ping(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("dasdasdad"))
}

func functionPointer(fp func(http.ResponseWriter, *http.Request)) (func(http.ResponseWriter, *http.Request), error) {
	
}

func main() {
	http.HandleFunc("/", ping)
	http.ListenAndServe(":8090", nil)
}
