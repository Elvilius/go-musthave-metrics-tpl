package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/app/agent"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	agent, err := agent.New()
	if err != nil {
		slog.Error("Error start app")
	}
	agent.Run(ctx)
}
