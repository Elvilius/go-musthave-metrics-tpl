package models

import "encoding/json"

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metrics) MarshalValue() ([]byte, error) {
	if m.MType == Counter {
		return json.Marshal(m.Delta)
	} else {
		return json.Marshal(m.Value)
	}
}
