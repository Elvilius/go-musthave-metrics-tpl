package collector

import (
	"testing"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_GetMetrics(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	sugarLogger := logger.Sugar()

	testCfg := config.AgentConfig{PollInterval: 10, ServerAddress: "localhost:8080", ReportInterval: 3}
	agentServiceMetrics := New(&testCfg, sugarLogger)

	agentServiceMetrics.CollectMetric()
	metrics := agentServiceMetrics.GetMetrics()

	expectedGauges := map[string]struct{}{
		"Alloc":           {},
		"BuckHashSys":     {},
		"Frees":           {},
		"GCCPUFraction":   {},
		"GCSys":           {},
		"HeapAlloc":       {},
		"HeapIdle":        {},
		"HeapInuse":       {},
		"HeapObjects":     {},
		"HeapReleased":    {},
		"HeapSys":         {},
		"LastGC":          {},
		"Lookups":         {},
		"MCacheInuse":     {},
		"MCacheSys":       {},
		"MSpanInuse":      {},
		"MSpanSys":        {},
		"Mallocs":         {},
		"NextGC":          {},
		"NumForcedGC":     {},
		"NumGC":           {},
		"OtherSys":        {},
		"PauseTotalNs":    {},
		"StackInuse":      {},
		"StackSys":        {},
		"Sys":             {},
		"TotalAlloc":      {},
		"RandomValue":     {},
		"PollCount":       {},
		"TotalMemory":     {},
		"FreeMemory":      {},
		"CPUutilization1": {},
	}

	t.Run("expected Metrics", func(t *testing.T) {
		for _, m := range metrics {
			_, ok := expectedGauges[m.ID]
			assert.True(t, ok, m.ID)
		}
	})
}
