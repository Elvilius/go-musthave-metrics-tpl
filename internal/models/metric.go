package models

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type IMetric interface {
	Update()
	Get() Metrics
}

type AllocMetric struct {
	metric *Metrics
}

type BuckHashSysMetric struct {
	metric *Metrics
}

type FreesMetric struct {
	metric *Metrics
}

type GCCPUFractionMetric struct {
	metric *Metrics
}

type GCSysMetric struct {
	metric *Metrics
}

type HeapAllocMetric struct {
	metric *Metrics
}

type HeapIdleMetric struct {
	metric *Metrics
}

type HeapInuseMetric struct {
	metric *Metrics
}

type HeapObjectsMetric struct {
	metric *Metrics
}

type HeapReleasedMetric struct {
	metric *Metrics
}

type HeapSysMetric struct {
	metric *Metrics
}

type LastGCMetric struct {
	metric *Metrics
}

type LookupsMetric struct {
	metric *Metrics
}

type MCacheInuseMetric struct {
	metric *Metrics
}

type MCacheSysMetric struct {
	metric *Metrics
}

type MSpanInuseMetric struct {
	metric *Metrics
}

type MSpanSysMetric struct {
	metric *Metrics
}

type MallocsMetric struct {
	metric *Metrics
}

type NextGCMetric struct {
	metric *Metrics
}

type NumForcedGCMetric struct {
	metric *Metrics
}

type NumGCMetric struct {
	metric *Metrics
}

type OtherSysMetric struct {
	metric *Metrics
}

type PauseTotalNsMetric struct {
	metric *Metrics
}

type StackInuseMetric struct {
	metric *Metrics
}

type StackSysMetric struct {
	metric *Metrics
}

type SysMetric struct {
	metric *Metrics
}

type TotalAllocMetric struct {
	metric *Metrics
}

type RandomValueMetric struct {
	metric *Metrics
}

type PollCountMetric struct {
	metric *Metrics
	count  int64
}

type TotalMemoryMetric struct {
	metric *Metrics
}

type FreeMemoryMetric struct {
	metric *Metrics
}

type CPUutilization1Metric struct {
	metric *Metrics
}
