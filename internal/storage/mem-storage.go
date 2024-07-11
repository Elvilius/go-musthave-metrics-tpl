package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
)

type MemStorage struct {
	metrics map[string]models.Metrics
	cfg     *config.ServerConfig
}

func NewMemStorage(cfg *config.ServerConfig) handler.Storager {
	return &MemStorage{metrics: make(map[string]models.Metrics), cfg: cfg}
}

func (r *MemStorage) Save(metric models.Metrics) error {
	mType, ID, value, delta := metric.MType, metric.ID, metric.Value, metric.Delta

	var defaultDelta int64 = 0
	if delta == nil {
		delta = &defaultDelta
	}
	existMetric, ok := r.Get(mType, ID)

	if mType == models.Gauge {
		r.metrics[ID] = models.Metrics{ID: ID, MType: mType, Value: value}
		return nil
	}
	if mType == models.Counter {
		if !ok {
			r.metrics[ID] = models.Metrics{ID: ID, MType: mType, Delta: delta}
			return nil
		} else {
			delta := *existMetric.Delta + *delta
			existMetric.Delta = &delta
			r.metrics[ID] = existMetric
		}
	}
	return nil
}

func (r *MemStorage) Get(mType string, ID string) (models.Metrics, bool) {
	m, ok := r.metrics[ID]
	if !ok {
		return models.Metrics{}, false
	}
	if m.MType != mType {
		return models.Metrics{}, false
	}

	return m, true
}

func (r *MemStorage) GetAll() []models.Metrics {
	all := make([]models.Metrics, 0, len(r.metrics))
	for _, m := range r.metrics {
		all = append(all, m)
	}
	return all
}

func (r *MemStorage) SaveToFile() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(wd, r.cfg.FileStoragePath)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	metrics := r.GetAll()
	bytes, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (r *MemStorage) LoadFromFile() error {
	if r.cfg.FileStoragePath == "" || !r.cfg.Restore {
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(wd, r.cfg.FileStoragePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var loadedMetrics []models.Metrics

	if err := json.Unmarshal(data, &loadedMetrics); err != nil {
		return err
	}

	for _, metric := range loadedMetrics {
		err := r.Save(metric)
		if err != nil {
			return err
		}
	}

	return nil
}
