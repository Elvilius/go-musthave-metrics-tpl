package models

import "encoding/json"

const (
	Gauge   = "gauge"
	Counter = "counter"
)

const (
	MetricAlloc            = "Alloc"
	MetricBuckHashSys      = "BuckHashSys"
	MetricFrees            = "Frees"
	MetricGCCPUFraction    = "GCCPUFraction"
	MetricGCSys            = "GCSys"
	MetricHeapAlloc        = "HeapAlloc"
	MetricHeapIdle         = "HeapIdle"
	MetricHeapInuse        = "HeapInuse"
	MetricHeapObjects      = "HeapObjects"
	MetricHeapReleased     = "HeapReleased"
	MetricHeapSys          = "HeapSys"
	MetricLastGC           = "LastGC"
	MetricLookups          = "Lookups"
	MetricMCacheInuse      = "MCacheInuse"
	MetricMCacheSys        = "MCacheSys"
	MetricMSpanInuse       = "MSpanInuse"
	MetricMSpanSys         = "MSpanSys"
	MetricMallocs          = "Mallocs"
	MetricNextGC           = "NextGC"
	MetricNumForcedGC      = "NumForcedGC"
	MetricNumGC            = "NumGC"
	MetricOtherSys         = "OtherSys"
	MetricPauseTotalNs     = "PauseTotalNs"
	MetricStackInuse       = "StackInuse"
	MetricStackSys         = "StackSys"
	MetricSys              = "Sys"
	MetricTotalAlloc       = "TotalAlloc"
	MetricRandomValue      = "RandomValue"
	MetricPollCount        = "PollCount"
	MetricTotalMemory      = "TotalMemory"
	MetricFreeMemory       = "FreeMemory"
	MetricCPUUtilization1  = "CPUutilization1"
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

func (m *Metrics) MarshalMetric() ([]byte, error) {
	return json.Marshal(m)
}
