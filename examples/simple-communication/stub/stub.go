package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	bhttp "github.com/alexcogojocaru/btracer/http"
	"github.com/alexcogojocaru/btracer/trace"
)

func StubRoute(w http.ResponseWriter, req *http.Request) {
	provider := req.Context().Value("provider").(*trace.TraceProvider)
	ctx, span := provider.Start(req.Context(), "StubRoute")
	span.AddLog("ERROR", "Stub Route Context")
	span.End()

	ctxList := []context.Context{ctx}
	for idx := 0; idx < rand.Intn(5); idx++ {
		ctxIdx := rand.Intn(len(ctxList))
		localCtx, localSpan := provider.Start(ctxList[ctxIdx], fmt.Sprintf("Stub-Loop-%d", idx))
		defer localSpan.End()

		localSpan.AddLog("INFO", fmt.Sprintf("Testing loop %d", idx))
		ctxList = append(ctxList, localCtx)
	}

	w.Write([]byte("StubRoute is working"))
}

func main() {
	http.Handle("/stub", bhttp.NewHandler(StubRoute, "Stub", "Stub_Microservice"))
	http.ListenAndServe("0.0.0.0:8091", nil)
}
