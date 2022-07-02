package http

import (
	"context"
	"io"
	"net/http"

	"github.com/alexcogojocaru/btracer/propagation"
)

type TraceClient struct {
	Client     http.Client
	Propagator propagation.Propagator
}

func NewTraceClient() *TraceClient {
	return &TraceClient{
		Client:     http.Client{},
		Propagator: propagation.Propagator{},
	}
}

func (tc *TraceClient) Request(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	tc.Propagator.Inject(ctx, req)
	return tc.Client.Do(req)
}
