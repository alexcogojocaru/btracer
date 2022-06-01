package main

import (
	"context"
	"time"

	"github.com/alexcogojocaru/btracer/trace"
)

func main() {
	provider := trace.NewProvider()
	defer provider.Shutdown()

	for {
		ctx1, _ := provider.Start(context.Background(), "Main")
		provider.Start(ctx1, "SecondMain")
		ctx3, _ := provider.Start(ctx1, "ThirdMain")
		provider.Start(ctx3, "FourthMain")

		time.Sleep(time.Second * 3)
	}
}
