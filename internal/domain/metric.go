package domain

type MetricType string

type Metric struct {
	Name  string
	Type  MetricType
	Value any
}

const (
	Gauge   = "gauge"
	Counter = "counter"
)
