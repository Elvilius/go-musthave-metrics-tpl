package storage

import "github.com/Elvilius/go-musthave-metrics-tpl/internal/domain"

type Storage interface {
	Save(metricType string, metricName string, value any) error
	Get(metricType, metricName string) (domain.Metric, bool)
	GetAll() []domain.Metric
}
