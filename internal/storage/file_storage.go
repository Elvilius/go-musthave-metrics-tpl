package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/metrics"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
)

type FileStorage struct {
	cfg     *config.ServerConfig
	storage metrics.Storager
}

func NewFileStorage(cfg *config.ServerConfig, storage metrics.Storager) *FileStorage {
	return &FileStorage{cfg: cfg, storage: storage}
}

func (f *FileStorage) SaveToFile() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(wd, f.cfg.FileStoragePath)
	tempPath := path + ".tmp"

	file, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	metrics, err := f.storage.GetAll(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to get metrics: %w", err)
	}

	bytes, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	err = os.Rename(tempPath, path)
	if err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}
func (f *FileStorage) LoadFromFile() error {
	if f.cfg.FileStoragePath == "" || !f.cfg.Restore {
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(wd, f.cfg.FileStoragePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var loadedMetrics []models.Metrics

	if err := json.Unmarshal(data, &loadedMetrics); err != nil {
		return err
	}

	for _, metric := range loadedMetrics {
		err := f.storage.Save(context.TODO(), metric)
		if err != nil {
			return err
		}
	}

	return nil
}
