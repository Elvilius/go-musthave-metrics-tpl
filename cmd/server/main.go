package main

import (
	"database/sql"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/server"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/storage"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/logger"
	_ "github.com/lib/pq"
)

func main() {
	logger, err := logger.New()
	if err != nil {
		panic(err)
	}
	cfg := config.NewServer()

	db, err := sql.Open("postgres", cfg.DatabaseDsn)
	if err != nil {
		logger.Fatalw("Failed to open DB", "error", err)
	}
	defer db.Close()

	storage := storage.New(cfg, db, logger)
	handler := handler.NewHandler(storage)

	server := server.New(cfg, handler, logger, db)

	server.Run()
}
