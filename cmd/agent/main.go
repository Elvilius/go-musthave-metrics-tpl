package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/app/agent"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/logger"
)

var (
	BuildVersion string = "NA"
	BuildDate    string = "NA"
	BuildCommit  string = "NA"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, err := logger.New()
	if err != nil {
		logger.Fatal(err)
	}

	
	agent, err := agent.New(logger)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Build version=%s \n", BuildVersion)
	logger.Infof("Build date=%s \n", BuildDate)
	logger.Infof("Build commit=%s \n", BuildCommit)

	agent.Run(ctx)
	os.Exit(1)
}
