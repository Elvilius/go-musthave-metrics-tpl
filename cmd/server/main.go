package main

import (
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/server"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/storage"
)

func main() {
	memStorage := storage.NewMemStorage()
	handler := handler.NewHandler(memStorage)
	cfg := config.GetServerConfig()
	server := server.NewServer(&cfg, handler)

	server.Run()
}
