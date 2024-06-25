package main

import (
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/services"
)

func main() {
	cfg := config.GetAgentConfig()
	agent := services.NewAgentMetricService(cfg)
	agent.Run()
}
