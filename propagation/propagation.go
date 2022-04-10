package propagation

import (
	"context"

	bagent "github.com/alexcogojocaru/btracer/proto-gen/btrace_agent"
)

type Carrier struct {
	Values map[string]string
}

type Propagator interface {
	Inject(ctx context.Context, span bagent.Span) error
	Extract(ctx context.Context) error
}

type propagator struct {
}

func (p *propagator) Inject(ctx context.Context, span bagent.Span) {

}
