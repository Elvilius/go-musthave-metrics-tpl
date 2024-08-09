package storage

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func New(cfg *config.ServerConfig, db *sql.DB, logger *zap.SugaredLogger) handler.Storager {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.DatabaseDsn == "" {
		memStorage := NewMemStorage(cfg)
		fs := NewFileStorage(cfg, memStorage)
		go runFile(ctx, cfg, fs, logger)
		return memStorage
	}

	if err := goose.Up(db, "./internal/storage/migrations"); err != nil {
		logger.Fatalw("Failed to run migrations", "error", err)
	}
	return NewDBStorage(db)
}

func runFile(ctx context.Context, cfg *config.ServerConfig, fs *FileStorage, logger *zap.SugaredLogger) {
	ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
	defer ticker.Stop()

	wd, err := os.Getwd()
	if err != nil {
		logger.Errorln("Failed to get working directory:", err)
	}
	dir, _ := filepath.Split(cfg.FileStoragePath)
	if err := os.MkdirAll(filepath.Join(wd, dir), 0o777); err != nil {
		logger.Errorln("Failed to create directories:", err)
	}

	err = fs.LoadFromFile()
	if err != nil {
		logger.Errorln("Failed to load from file:", err)
	}

	for {
		select {
		case <-ticker.C:
			err := fs.SaveToFile()
			if err != nil {
				logger.Errorln("Failed to save to file:", err)
			}
		case <-ctx.Done():
			err := fs.SaveToFile()
			if err != nil {
				logger.Errorln("Failed to save to file during shutdown:", err)
			}
			return
		}
	}
}
