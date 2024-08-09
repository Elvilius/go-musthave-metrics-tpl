package storage

import (
	"context"

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

func (r *MemStorage) Save(ctx context.Context, metric models.Metrics) error {
	mType, ID, value, delta := metric.MType, metric.ID, metric.Value, metric.Delta

	var defaultDelta int64 = 0
	if delta == nil {
		delta = &defaultDelta
	}
	existMetric, ok, err := r.Get(ctx, mType, ID)
	if err != nil {
		return err
	}

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

func (r *MemStorage) Get(ctx context.Context, mType string, ID string) (models.Metrics, bool, error) {
	m, ok := r.metrics[ID]
	if !ok {
		return models.Metrics{}, false, nil
	}
	if m.MType != mType {
		return models.Metrics{}, false, nil
	}

	return m, true, nil
}

func (r *MemStorage) GetAll(ctx context.Context) ([]models.Metrics, error) {
	all := make([]models.Metrics, 0, len(r.metrics))
	for _, m := range r.metrics {
		all = append(all, m)
	}
	return all, nil
}
