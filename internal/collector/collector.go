package collector

import (
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"go.uber.org/zap"
)

type Collector struct {
	cfg       *config.AgentConfig
	metrics   []models.IMetric
	pollCount int64
	logger    *zap.SugaredLogger
}

func New(cfg *config.AgentConfig, logger *zap.SugaredLogger) *Collector {
	metrics := []models.IMetric{
		models.NewAllocMetric(),
		models.NewBuckHashSysMetric(),
		models.NewFreesMetric(),
		models.NewGCCPUFractionMetric(),
		models.NewGCSysMetric(),
		models.NewHeapAllocMetric(),
		models.NewHeapIdleMetric(),
		models.NewHeapInuseMetric(),
		models.NewHeapObjectsMetric(),
		models.NewHeapReleasedMetric(),
		models.NewHeapSysMetric(),
		models.NewLastGCMetric(),
		models.NewLookupsMetric(),
		models.NewMCacheInuseMetric(),
		models.NewMSpanInuseMetric(),
		models.NewMSpanSysMetric(),
		models.NewMallocsMetric(),
		models.NewNextGCMetric(),
		models.NewNumForcedGC(),
		models.NewNumGCMetric(),
		models.NewOtherSysMetric(),
		models.NewPauseTotalNsMetric(),
		models.NewRandomValueMetric(),
		models.NewPollCountMetric(),
		models.NewTotalAlloc(),
		models.NewMCacheSys(),
		models.NewStackInuseMetric(),
		models.NewStackSysMetric(),
		models.NewSysMetric(),
		models.NewFreeMemoryMetric(),
		models.NewCPUutilization1Metric(),
		models.NewTotalMemoryMetric(),
	}

	return &Collector{cfg: cfg, metrics: metrics, pollCount: 0, logger: logger}
}

func (s *Collector) CollectMetric() {
	for _, m := range s.metrics {
		m.Update()
	}
}

func (s *Collector) GetMetrics() []models.Metrics {
	metrics := make([]models.Metrics, len(s.metrics))
	for i, metric := range s.metrics {
		metrics[i] = metric.Get()
	}
	return metrics
}
