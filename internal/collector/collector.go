package collector

import (
	"math/rand/v2"
	"runtime"
	"sync/atomic"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/shirou/gopsutil/v4/cpu"
	"go.uber.org/zap"

	"github.com/shirou/gopsutil/v4/mem"
)

var memStats runtime.MemStats

type Collector struct {
	cfg     *config.AgentConfig
	metrics []*models.Metrics
	logger  *zap.SugaredLogger
}

func New(cfg *config.AgentConfig, logger *zap.SugaredLogger) *Collector {
	createDelta := func(n int64) *int64 {
		return &n
	}
	metrics := []*models.Metrics{
		{ID: models.MetricAlloc, MType: models.Gauge},
		{ID: models.MetricBuckHashSys, MType: models.Gauge},
		{ID: models.MetricFrees, MType: models.Gauge},
		{ID: models.MetricGCCPUFraction, MType: models.Gauge},
		{ID: models.MetricGCSys, MType: models.Gauge},
		{ID: models.MetricHeapAlloc, MType: models.Gauge},
		{ID: models.MetricHeapIdle, MType: models.Gauge},
		{ID: models.MetricHeapInuse, MType: models.Gauge},
		{ID: models.MetricHeapObjects, MType: models.Gauge},
		{ID: models.MetricHeapReleased, MType: models.Gauge},
		{ID: models.MetricHeapSys, MType: models.Gauge},
		{ID: models.MetricLastGC, MType: models.Gauge},
		{ID: models.MetricLookups, MType: models.Gauge},
		{ID: models.MetricMCacheInuse, MType: models.Gauge},
		{ID: models.MetricMCacheSys, MType: models.Gauge},
		{ID: models.MetricMSpanInuse, MType: models.Gauge},
		{ID: models.MetricMSpanSys, MType: models.Gauge},
		{ID: models.MetricMallocs, MType: models.Gauge},
		{ID: models.MetricNextGC, MType: models.Gauge},
		{ID: models.MetricNumForcedGC, MType: models.Gauge},
		{ID: models.MetricNumGC, MType: models.Gauge},
		{ID: models.MetricOtherSys, MType: models.Gauge},
		{ID: models.MetricPauseTotalNs, MType: models.Gauge},
		{ID: models.MetricStackInuse, MType: models.Gauge},
		{ID: models.MetricStackSys, MType: models.Gauge},
		{ID: models.MetricTotalAlloc, MType: models.Gauge},
		{ID: models.MetricRandomValue, MType: models.Gauge},
		{ID: models.MetricPollCount, MType: models.Counter, Delta: createDelta(0)},
		{ID: models.MetricTotalMemory, MType: models.Gauge},
		{ID: models.MetricFreeMemory, MType: models.Gauge},
		{ID: models.MetricCPUUtilization1, MType: models.Gauge},
	}
	return &Collector{cfg: cfg, metrics: metrics, logger: logger}
}

func (c *Collector) CollectMetric() {
	runtime.ReadMemStats(&memStats)
	memInfo, _ := mem.VirtualMemory()
	cpuInfo, _ := cpu.Percent(0, true)

	for _, m := range c.metrics {

		var value float64

		switch m.ID {
		case models.MetricPollCount:
			atomic.AddInt64(m.Delta, 1)
		case models.MetricRandomValue:
			value = rand.Float64()
		case models.MetricCPUUtilization1:
			value = cpuInfo[1]
		case models.MetricTotalMemory:
			value = float64(memInfo.Total)
		case models.MetricFreeMemory:
			memInfo, _ := mem.VirtualMemory()
			value = float64(memInfo.Free)
		case models.MetricAlloc:
			value = float64(memStats.Alloc)
		case models.MetricBuckHashSys:
			value = float64(memStats.BuckHashSys)
		case models.MetricFrees:
			value = float64(memStats.Frees)
		case models.MetricGCCPUFraction:
			value = float64(memStats.GCCPUFraction)
		case models.MetricGCSys:
			value = float64(memStats.GCSys)
		case models.MetricHeapAlloc:
			value = float64(memStats.HeapAlloc)
		case models.MetricHeapIdle:
			value = float64(memStats.HeapIdle)
		case models.MetricHeapInuse:
			value = float64(memStats.HeapInuse)
		case models.MetricHeapObjects:
			value = float64(memStats.HeapObjects)
		case models.MetricHeapReleased:
			value = float64(memStats.HeapReleased)
		case models.MetricHeapSys:
			value = float64(memStats.HeapSys)
		case models.MetricLastGC:
			value = float64(memStats.LastGC)
		case models.MetricLookups:
			value = float64(memStats.Lookups)
		case models.MetricMCacheInuse:
			value = float64(memStats.MCacheInuse)
		case models.MetricMCacheSys:
			value = float64(memStats.MCacheSys)
		case models.MetricMSpanInuse:
			value = float64(memStats.MSpanInuse)
		case models.MetricMSpanSys:
			value = float64(memStats.MSpanSys)
		case models.MetricMallocs:
			value = float64(memStats.Mallocs)
		case models.MetricNextGC:
			value = float64(memStats.NextGC)
		case models.MetricNumForcedGC:
			value = float64(memStats.NumForcedGC)
		case models.MetricNumGC:
			value = float64(memStats.NumGC)
		case models.MetricOtherSys:
			value = float64(memStats.OtherSys)
		case models.MetricPauseTotalNs:
			value = float64(memStats.PauseTotalNs)
		case models.MetricStackInuse:
			value = float64(memStats.StackInuse)
		case models.MetricStackSys:
			value = float64(memStats.StackSys)
		case models.MetricSys:
			value = float64(memStats.Sys)
		case models.MetricTotalAlloc:
			value = float64(memStats.TotalAlloc)
		default:
			c.logger.Errorln("uknow metric ID")
		}
		m.Value = &value
	}
}

func (c *Collector) GetMetrics() []*models.Metrics {
	return c.metrics
}
