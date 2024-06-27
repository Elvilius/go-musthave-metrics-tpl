package main

import (
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/server"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/storage"
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

	memStorage := storage.NewMemStorage()
	handler := handler.NewHandler(memStorage)
	cfg := config.GetServerConfig()
	server := server.New(&cfg, handler, sugarLogger)

	server.Run()
}
