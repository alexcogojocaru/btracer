package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/alexcogojocaru/btracer/config"
	thttp "github.com/alexcogojocaru/btracer/http"
	"github.com/alexcogojocaru/btracer/trace"
)

const LISTENER_URL = "http://localhost:8090/ping"

func main() {
	conf, err := config.ParseConfig("../../../config/config.yml")
	if err != nil {
		log.Fatal("Error while parsing the config file")
	}

	client := thttp.NewTraceClient()

	var exporterConfig trace.ExporterConfig
	if conf.Agent.Bypass == false {
		exporterConfig = trace.ExporterConfig{
			Bypass: false,
			AgentConfig: trace.ConnectionDetails{
				Host: conf.Agent.Hostname,
				Port: int(conf.Agent.Port),
			},
		}
	} else {
		exporterConfig = trace.ExporterConfig{
			Bypass: true,
			CollectorConfig: trace.ConnectionDetails{
				Host: conf.Collector.Hostname,
				Port: int(conf.Collector.Port),
			},
		}
	}

	provider := trace.NewProvider("Caller_Microservice", exporterConfig)
	defer provider.Shutdown()

	ctx, span := provider.Start(context.Background(), "Caller_Main")
	defer span.End()
	span.AddLog("INFO", "Starting a main block")

	client.Request(ctx, "GET", LISTENER_URL, nil)

	_, cacheSpan := provider.Start(ctx, "CacheCall")
	defer cacheSpan.End()
	cacheSpan.AddLog("ERROR", "redis cache call after the request")
	cacheSpan.AddLog("INFO", "something bad happened")

	dbCtx, dbSpan := provider.Start(ctx, "DbCall")
	defer dbSpan.End()
	dbSpan.AddLog("WARNING", "Connection to db is slow")

	ctxTemp, tempSpan := provider.Start(dbCtx, "TempCall")
	defer tempSpan.End()
	time.Sleep(3)
	_, tempSpan0 := provider.Start(ctxTemp, "TempCall")
	defer tempSpan0.End()

	contextList := []context.Context{ctx, dbCtx, ctxTemp}
	for idx := 0; idx < 10; idx++ {
		ctxIdx := rand.Intn(len(contextList))

		localContext, localSpan := provider.Start(contextList[ctxIdx], fmt.Sprintf("InsideLoop-%d", idx))
		defer localSpan.End()

		for logIdx := 0; logIdx < rand.Intn(4); logIdx++ {
			localSpan.AddLog("INFO", fmt.Sprintf("Log generated inside loop %d-%d", idx, logIdx))
		}

		contextList = append(contextList, localContext)
	}
}
