package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_collectMetrics(t *testing.T) {
	metrics := collectMetrics()
	expectedGauges := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
		"StackSys", "Sys", "TotalAlloc", "RandomValue",
	}

	expectedCounters := []string{"PollCount"}

	t.Run("expected Gauges", func(t *testing.T) {
		for _, expectedGauge := range expectedGauges {
			_, ok := metrics[expectedGauge]
			assert.True(t, ok)
		}
	})

	t.Run("expected Counters", func(t *testing.T) {
		for _, expectedCounter := range expectedCounters {
			_, ok := metrics[expectedCounter]
			assert.True(t, ok)
		}
	})
}
