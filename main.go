package main

func main() {
	// agentConfig := bee.AgentConfig{Host: "localhost", Port: 4576}
	// beeExporter, _ := bee.NewBeeExporter(&agentConfig)

	// traceProvider := trace.NewTracerProvider(trace.WithBatcher(beeExporter))
	// defer func() {
	// 	if err := traceProvider.Shutdown(context.Background()); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	// otel.SetTracerProvider(traceProvider)
	// otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
	// 	propagation.TraceContext{},
	// 	propagation.Baggage{},
	// ))

	// tracer := otel.Tracer("BTracer")

	// otelCtx, span := tracer.Start(context.Background(), "Main")
	// defer span.End()

	// _, span1 := tracer.Start(otelCtx, "SubMain")
	// time.Sleep(1 * time.Second)
	// defer span1.End()
}
