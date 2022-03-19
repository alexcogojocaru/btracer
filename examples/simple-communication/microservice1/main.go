package main

import (
	bee "github.com/alexcogojocaru/btracer/exporters/bee"
)

func main() {
	agentConfig := bee.AgentConfig{Host: "localhost", Port: 4576}
	beeExporter, _ := bee.NewBeeExporter(&agentConfig)
}
