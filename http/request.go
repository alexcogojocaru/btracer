package http

import (
	"context"
	"log"
	"net/http"
)

func Request(url string, ctx context.Context) (*http.Response, error) {
	log.Print(ctx)
	resp, err := http.Get(url)
	return resp, err
}
