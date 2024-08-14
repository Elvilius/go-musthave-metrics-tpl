package storage

import (
	"context"
	"sync"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
)

type MemStorage struct {
	metrics map[string]models.Metrics
	cfg     *config.ServerConfig
	rw      sync.RWMutex
}

func NewMemStorage(cfg *config.ServerConfig) handler.Storager {
	return &MemStorage{metrics: make(map[string]models.Metrics), cfg: cfg, rw: sync.RWMutex{}}
}

func (m *MemStorage) Save(ctx context.Context, metric models.Metrics) error {
	mType, ID, value, delta := metric.MType, metric.ID, metric.Value, metric.Delta

	var defaultDelta int64 = 0
	if delta == nil {
		delta = &defaultDelta
	}
	existMetric, ok, err := m.Get(ctx, mType, ID)
	if err != nil {
		return err
	}

	if mType == models.Gauge {
		m.rw.Lock()
		defer m.rw.Unlock()
		m.metrics[ID] = models.Metrics{ID: ID, MType: mType, Value: value}
		return nil

	}
	if mType == models.Counter {
		m.rw.Lock()
		defer m.rw.Unlock()
		if !ok {
			m.metrics[ID] = models.Metrics{ID: ID, MType: mType, Delta: delta}
			return nil
		} else {
			delta := *existMetric.Delta + *delta
			existMetric.Delta = &delta
			m.metrics[ID] = existMetric
		}
	}
	return nil
}

func (m *MemStorage) Get(ctx context.Context, mType, ID string) (models.Metrics, bool, error) {
	metric, ok := m.metrics[ID]
	if !ok {
		return models.Metrics{}, false, nil
	}
	if metric.MType != mType {
		return models.Metrics{}, false, nil
	}

	return metric, true, nil
}

func (m *MemStorage) GetAll(ctx context.Context) ([]models.Metrics, error) {
	all := make([]models.Metrics, 0, len(m.metrics))
	for _, m := range m.metrics {
		all = append(all, m)
	}
	return all, nil
}

func (m *MemStorage) Updates(ctx context.Context, metrics []models.Metrics) error {
	for _, metric := range metrics {
		err := m.Save(ctx, metric)
		if err != nil {
			return err
		}
	}
	return nil
}
