package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/app/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	app := server.New()
	app.Run(ctx)
}
