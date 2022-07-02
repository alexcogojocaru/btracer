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
	span.AddLog("INFO", "buna siua ce faceti")
	span.AddLog("INFO", "buna siua ce faceti 2")
	defer span.End()

	ctx1, span1 := provider.Start(ctx, "ThirdMain")
	defer span1.End()

	_, span2 := provider.Start(ctx1, "FourthMain")
	defer span2.End()
	_, span3 := provider.Start(ctx, "FifthMain")
	defer span3.End()
}
