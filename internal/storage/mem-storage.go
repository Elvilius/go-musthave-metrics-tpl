package storage

import (
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/domain"
)

type MemStorage struct {
	metrics map[string]domain.Metric
}

func NewMemStorage() Storage {
	return &MemStorage{metrics: make(map[string]domain.Metric)}
}

func (r *MemStorage) Save(metricType string, metricName string, value any) {
	existMetric, ok := r.Get(metricType, metricName)
	switch metricType {
	case domain.Gauge:
		r.metrics[metricName] = domain.Metric{Type: metricType, Name: metricName, Value: value}
	case domain.Counter:
		var newValue int
		if ok {
			newValue = existMetric.Value.(int) + value.(int)
		} else {
			newValue = value.(int)
		}
		r.metrics[metricName] = domain.Metric{Type: metricType, Name: metricName, Value: newValue}
	}
}

func (r *MemStorage) Get(metricType string, metricName string) (domain.Metric, bool) {
	m, ok := r.metrics[metricName]
	if !ok {
		return domain.Metric{}, false
	}
	if m.Type != metricType {
		return domain.Metric{}, false
	}

	return m, true
}

func (r *MemStorage) GetAll() []domain.Metric {
	all := make([]domain.Metric, 0, len(r.metrics))
	for _, m := range r.metrics {
		all = append(all, m)
	}
	return all
}
