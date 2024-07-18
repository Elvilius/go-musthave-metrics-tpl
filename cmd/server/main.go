package main

import (
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/server"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/storage"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/logger"
)

func main() {
	logger, err := logger.New()
	if err != nil {
		panic(err)
	}
	cfg := config.NewServer()
	memStorage := storage.NewMemStorage(cfg)
	handler := handler.NewHandler(memStorage)
	server := server.New(cfg, handler, logger)

	server.Run(memStorage)
}
