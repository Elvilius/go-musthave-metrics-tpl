package storage

import (
	"fmt"

	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
)

type MemStorage struct {
	metrics map[string]models.Metrics
}

func NewMemStorage() handler.Storager {
	return &MemStorage{metrics: make(map[string]models.Metrics)}
}

func (r *MemStorage) Save(metric models.Metrics) error {
	mType, ID, value, delta := metric.MType, metric.ID, metric.Value, metric.Delta

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
		fmt.Println(m)
		all = append(all, m)
	}
	return all
}
