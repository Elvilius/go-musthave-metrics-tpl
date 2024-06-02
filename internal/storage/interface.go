package storage

type Storage interface {
	Gauge(metricName string, value float64)
	Inc(metricName string)
}
