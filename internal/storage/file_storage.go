package storage

import (
	"bufio"
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
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	path := filepath.Join(wd, f.cfg.FileStoragePath)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
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

	writer := bufio.NewWriter(file)
	_, err = writer.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush data to file: %w", err)
	}

	if err := file.Truncate(int64(len(bytes))); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
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
