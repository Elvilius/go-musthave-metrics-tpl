package storage

import (
	"database/sql"
	"fmt"
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
	if cfg.DatabaseDsn == "" {
		memStorage := NewMemStorage(cfg)
		fs := NewFileStorage(cfg, memStorage)
		runFile(cfg, fs, logger)
		return memStorage
	}

	if err := goose.Up(db, "./internal/storage/migrations"); err != nil {
		logger.Fatalw("Failed to run migrations", "error", err)
	}
	return NewDbStorage(db)
}

func runFile(cfg *config.ServerConfig, fs *FileStorage, logger *zap.SugaredLogger) {
	ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
	defer ticker.Stop()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	dir, _ := filepath.Split(cfg.FileStoragePath)
	if err := os.MkdirAll(filepath.Join(wd, dir), 0o777); err != nil {
		fmt.Println(err)
	}

	err = fs.LoadFromFile()
	if err != nil {
		logger.Errorln(err)
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				err := fs.SaveToFile()
				if err != nil {
					logger.Errorln(err)
				}
			case <-done:
				err := fs.SaveToFile()
				if err != nil {
					logger.Errorln(err)
				}
				close(done)
				return
			}
		}
	}()

	<-done
}
