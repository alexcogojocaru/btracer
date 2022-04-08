package http

import (
	"context"
	"log"
	"net/http"
)

func Request(ctx context.Context, url string) (*http.Response, error) {
	log.Print(ctx)
	resp, err := http.Get(url)
	return resp, err
}
