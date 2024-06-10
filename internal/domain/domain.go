package domain

type Metric struct {
	Name  string
	Type  string
	Value any
}

const (
	Gauge   = "gauge"
	Counter = "counter"
)
