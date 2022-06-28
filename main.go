package main

import (
	"context"
	"time"

	"github.com/alexcogojocaru/btracer/trace"
)

func main() {
	provider := trace.NewProvider("TestingMain")
	defer provider.Shutdown()

	ctx, span := provider.Start(context.Background(), "Main")
	time.Sleep(5 * time.Millisecond)
	defer span.End()

	ctx3, span3 := provider.Start(ctx, "ThirdMain")
	defer span3.End()

	ctx4, span4 := provider.Start(ctx3, "FourthMain")
	defer span4.End()
	_, span5 := provider.Start(ctx4, "FifthMain")
	defer span5.End()
}
