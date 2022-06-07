package main

import (
	"context"

	"github.com/alexcogojocaru/btracer/trace"
)

func main() {
	provider := trace.NewProvider("BTracer")
	defer provider.Shutdown()

	ctx, _ := provider.Start(context.Background(), "Main")
	provider.Start(ctx, "SecondMain")
}
