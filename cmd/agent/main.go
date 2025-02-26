package main

import (
	"context"
	"fmt"
	//"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/app/agent"
)

var (
	BuildVersion string = "NA"
	BuildDate    string = "NA"
	BuildCommit  string = "NA"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	agent, err := agent.New()
	if err != nil {
		slog.Error("Error start app")
	}
	
	fmt.Printf("Build version=%s \n", BuildVersion)
	fmt.Printf("Build date=%s \n", BuildDate)
	fmt.Printf("Build commit=%s \n", BuildCommit)

	agent.Run(ctx)
}
