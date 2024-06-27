package storage

import (
	"strconv"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/domain"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
)

type MemStorage struct {
	metrics map[string]domain.Metric
}

func NewMemStorage() handler.Storager {
	return &MemStorage{metrics: make(map[string]domain.Metric)}
}

func (r *MemStorage) Save(metricType string, metricName string, value any) error {
	existMetric, ok := r.Get(metricType, metricName)

	if metricType == domain.Gauge {
		parsedValueFloat, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return err
		}
		r.metrics[metricName] = domain.Metric{Type: domain.MetricType(metricType), Name: metricName, Value: parsedValueFloat}
		return nil
	}
	if metricType == domain.Counter {
		parsedValue, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return err
		}

		if !ok {
			r.metrics[metricName] = domain.Metric{Type: domain.MetricType(metricType), Name: metricName, Value: parsedValue}
			return nil
		} else {
			existMetric.Value = existMetric.Value.(int64) + parsedValue
			r.metrics[metricName] = existMetric
		}
	}
	return nil
}

func (r *MemStorage) Get(metricType string, metricName string) (domain.Metric, bool) {
	m, ok := r.metrics[metricName]
	if !ok {
		return domain.Metric{}, false
	}
	if m.Type != domain.MetricType(metricType) {
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
