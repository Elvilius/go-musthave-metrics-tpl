package main

import (
	"context"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/services"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	logger, err := logger.New()
	if err != nil {
		panic(err)
	}
	cfg := config.NewAgent()
	agent := services.NewAgentMetricService(cfg, logger)
	agent.Run(ctx)
}
