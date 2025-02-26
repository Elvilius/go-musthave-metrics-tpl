package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/app/server"
)

var (
	BuildVersion string = "NA"
	BuildDate    string = "NA"
	BuildCommit  string = "NA"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	app := server.New()

	fmt.Printf("Build version=%s \n", BuildVersion)
	fmt.Printf("Build date=%s \n", BuildDate)
	fmt.Printf("Build commit=%s \n", BuildCommit)

	app.Run(ctx)
}
