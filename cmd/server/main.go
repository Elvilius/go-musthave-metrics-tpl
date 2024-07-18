package main

import (
	"context"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/server"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/storage"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/logger"
	"github.com/jackc/pgx/v5"
)

func main() {
	logger, err := logger.New()
	if err != nil {
		panic(err)
	}
	cfg := config.NewServer()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.DatabaseDsn)
	if err != nil {
		logger.Error("Error connect to database")
	}
	defer conn.Close(context.Background())

	memStorage := storage.NewMemStorage(cfg)
	handler := handler.NewHandler(memStorage)
	server := server.New(cfg, handler, logger, conn)

	server.Run(memStorage)
}
