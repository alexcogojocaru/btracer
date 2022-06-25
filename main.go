package main

import (
	"context"
	"time"

	"github.com/alexcogojocaru/btracer/trace"
)

func main() {
	provider := trace.NewProvider("TestingMain")
	defer provider.Shutdown()

	_, span := provider.Start(context.Background(), "Main")
	time.Sleep(5 * time.Millisecond)
	defer span.End()
	// provider.Start(ctx, "SecondMain")
	// ctx3, _ := provider.Start(ctx, "ThirdMain")
	// ctx4, _ := provider.Start(ctx3, "FourthMain")
	// provider.Start(ctx4, "FifthMain")
}
