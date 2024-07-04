package storage

import (
	"strconv"

	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
)

type MemStorage struct {
	metrics map[string]models.Metrics
}

func NewMemStorage() handler.Storager {
	return &MemStorage{metrics: make(map[string]models.Metrics)}
}

func (r *MemStorage) Save(metricType string, metricName string, value any) error {
	existMetric, ok := r.Get(metricType, metricName)

	if metricType == models.Gauge {
		parsedValueFloat, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return err
		}
		r.metrics[metricName] = models.Metrics{ID: metricName, MType: metricType, Value: &parsedValueFloat}
		return nil
	}
	if metricType == models.Counter {
		parsedValue, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return err
		}

		if !ok {
			r.metrics[metricName] = models.Metrics{ID: metricName, MType: metricType, Delta: &parsedValue}
			return nil
		} else {
			delta := *existMetric.Delta + parsedValue
			existMetric.Delta = &delta
			r.metrics[metricName] = existMetric
		}
	}
	return nil
}

func (r *MemStorage) Get(metricType string, metricName string) (models.Metrics, bool) {
	m, ok := r.metrics[metricName]
	if !ok {
		return models.Metrics{}, false
	}
	if m.MType != metricType {
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
