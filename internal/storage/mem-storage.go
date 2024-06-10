package storage

import (
	"fmt"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/domain"
)

type MemStorage struct {
	metrics map[string]domain.Metric
}

func NewMemStorage() Storage {
	return &MemStorage{metrics: make(map[string]domain.Metric)}
}

func (r *MemStorage) Save(metricType string, metricName string, value any) {
	if metricType == domain.Gauge {
		r.metrics[metricName] = domain.Metric{Type: metricType, Name: metricName, Value: value}
	} else if metricType == domain.Counter {
		existMetric, ok := r.Get(metricType, metricName)
		if !ok {
			r.metrics[metricName] = domain.Metric{Type: metricType, Name: metricName, Value: 1}
		} else {
			existMetric.Value = existMetric.Value.(int) + 1
			r.metrics[metricName] = existMetric
		}
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

func (r *MemStorage) Print() []domain.Metric {
	all := make([]domain.Metric, 0, len(r.metrics))
	for _, m := range r.metrics {
		fmt.Println(m)
	}
	return all
}
