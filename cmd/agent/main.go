package main

import (
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/services"
)

func main() {
	cfg := config.GetAgentConfig()
	agentServiceMetrics := services.NewAgentMetricService(cfg)
	agentServiceMetrics.SendMetrics()
}
