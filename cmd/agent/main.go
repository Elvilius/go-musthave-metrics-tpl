package main

import (
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/services"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	sugarLogger := logger.Sugar()

	cfg := config.GetAgentConfig()
	agent := services.NewAgentMetricService(cfg, sugarLogger)
	agent.Run()
}
