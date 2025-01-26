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
	return &MemStorage{
		metrics: make(map[string]models.Metrics),
		cfg:     cfg,
		rw:      sync.RWMutex{},
	}
}

func (m *MemStorage) Save(ctx context.Context, metric models.Metrics) error {
	mType, ID, value, delta := metric.MType, metric.ID, metric.Value, metric.Delta

	var defaultDelta int64 = 0
	if delta == nil {
		delta = &defaultDelta
	}

	m.rw.Lock()
	defer m.rw.Unlock()

	existMetric, ok := m.metrics[ID]

	if mType == models.Gauge {
		m.metrics[ID] = models.Metrics{ID: ID, MType: mType, Value: value}
		return nil
	}

	if mType == models.Counter {
		if !ok {
			m.metrics[ID] = models.Metrics{ID: ID, MType: mType, Delta: delta}
		} else {
			newDelta := *existMetric.Delta + *delta
			existMetric.Delta = &newDelta
			m.metrics[ID] = existMetric
		}
		return nil
	}

	return nil
}

func (m *MemStorage) Get(ctx context.Context, mType, ID string) (models.Metrics, bool, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()

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
	m.rw.RLock()
	defer m.rw.RUnlock()

	all := make([]models.Metrics, 0, len(m.metrics))
	for _, metric := range m.metrics {
		all = append(all, metric)
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
