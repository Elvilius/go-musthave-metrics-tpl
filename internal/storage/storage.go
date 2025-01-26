package storage

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/metrics"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type storageType string

const memType storageType = "1"
const dbType storageType = "2"

type Store struct {
	Storage metrics.Storager
	sType   storageType
	cfg     *config.ServerConfig
	db      *sql.DB
	logger  *zap.SugaredLogger
}

func New(cfg *config.ServerConfig, logger *zap.SugaredLogger, db *sql.DB) *Store {
	s := &Store{db: db, cfg: cfg, logger: logger}
	if cfg.DatabaseDsn == "" {
		s.sType = memType
		s.Storage = NewMemStorage(cfg)
	} else {
		s.sType = dbType
		s.Storage = NewDBStorage(db)
	}
	return s
}

func (s *Store) Run(ctx context.Context) {

	switch s.sType {
	case memType:
		fs := NewFileStorage(s.cfg, s.Storage)
		go s.runFile(ctx, s.cfg, fs, s.logger)
	case dbType:
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		if err := goose.UpContext(ctx, s.db, "./internal/storage/migrations"); err != nil {
			s.logger.Fatalw("Failed to run migrations", "error", err)
		}

	default:
		s.logger.Fatalw("Unknown store type")
	}
}

func (s *Store) Ping() error {
	if s.db == nil {
		return errors.New("DB not found")
	}

	return s.db.Ping()
}

func (s *Store) runFile(ctx context.Context, cfg *config.ServerConfig, fs *FileStorage, logger *zap.SugaredLogger) {
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
			os.Exit(1)
			return
		}
	}
}

func (s *Store) GetStorage() metrics.Storager {
	return s.Storage
}
