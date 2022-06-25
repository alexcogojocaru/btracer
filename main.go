package main

import (
	"context"

	"github.com/alexcogojocaru/btracer/trace"
)

func main() {
	provider := trace.NewProvider("TestingMain")
	defer provider.Shutdown()

	ctx, _ := provider.Start(context.Background(), "Main")
	provider.Start(ctx, "SecondMain")
	ctx3, _ := provider.Start(ctx, "ThirdMain")
	ctx4, _ := provider.Start(ctx3, "FourthMain")
	provider.Start(ctx4, "FifthMain")
}
