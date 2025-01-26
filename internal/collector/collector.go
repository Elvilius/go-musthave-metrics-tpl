package collector

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"go.uber.org/zap"
)

type Collector struct {
	cfg       *config.AgentConfig
	metrics   map[string]models.Metrics
	pollCount int64
	logger    *zap.SugaredLogger
}

func New(cfg *config.AgentConfig, logger *zap.SugaredLogger) *Collector {
	return &Collector{cfg: cfg, metrics: make(map[string]models.Metrics), pollCount: 0, logger: logger}
}

func (c *Collector) GetMetric() map[string]models.Metrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	alloc := float64(memStats.Alloc)
	c.metrics["Alloc"] = models.Metrics{ID: "Alloc", MType: models.Gauge, Value: &alloc}

	buckHashSys := float64(memStats.Frees)
	c.metrics["BuckHashSys"] = models.Metrics{ID: "BuckHashSys", MType: models.Gauge, Value: &buckHashSys}

	frees := float64(memStats.Frees)
	c.metrics["Frees"] = models.Metrics{ID: "Frees", MType: models.Gauge, Value: &frees}

	gCCPUFraction := memStats.GCCPUFraction
	c.metrics["GCCPUFraction"] = models.Metrics{ID: "GCCPUFraction", MType: models.Gauge, Value: &gCCPUFraction}

	gCSys := float64(memStats.GCSys)
	c.metrics["GCSys"] = models.Metrics{ID: "GCSys", MType: models.Gauge, Value: &gCSys}

	heapAlloc := float64(memStats.HeapAlloc)
	c.metrics["HeapAlloc"] = models.Metrics{ID: "HeapAlloc", MType: models.Gauge, Value: &heapAlloc}

	heapIdle := float64(memStats.HeapIdle)
	c.metrics["HeapIdle"] = models.Metrics{ID: "HeapIdle", MType: models.Gauge, Value: &heapIdle}

	heapInuse := float64(memStats.HeapInuse)
	c.metrics["HeapInuse"] = models.Metrics{ID: "HeapInuse", MType: models.Gauge, Value: &heapInuse}

	heapObjects := float64(memStats.HeapObjects)
	c.metrics["HeapObjects"] = models.Metrics{ID: "HeapObjects", MType: models.Gauge, Value: &heapObjects}

	heapReleased := float64(memStats.HeapReleased)
	c.metrics["HeapReleased"] = models.Metrics{ID: "HeapReleased", MType: models.Gauge, Value: &heapReleased}

	heapSys := float64(memStats.HeapSys)
	c.metrics["HeapSys"] = models.Metrics{ID: "HeapSys", MType: models.Gauge, Value: &heapSys}

	lastGC := float64(memStats.LastGC)
	c.metrics["LastGC"] = models.Metrics{ID: "LastGC", MType: models.Gauge, Value: &lastGC}

	lookups := float64(memStats.Lookups)
	c.metrics["Lookups"] = models.Metrics{ID: "Lookups", MType: models.Gauge, Value: &lookups}

	mCacheInuse := float64(memStats.MCacheInuse)
	c.metrics["MCacheInuse"] = models.Metrics{ID: "MCacheInuse", MType: models.Gauge, Value: &mCacheInuse}

	mCacheSys := float64(memStats.MCacheSys)
	c.metrics["MCacheSys"] = models.Metrics{ID: "MCacheSys", MType: models.Gauge, Value: &mCacheSys}

	mSpanInuse := float64(memStats.MSpanInuse)
	c.metrics["MSpanInuse"] = models.Metrics{ID: "MSpanInuse", MType: models.Gauge, Value: &mSpanInuse}

	mSpanSys := float64(memStats.MSpanSys)
	c.metrics["MSpanSys"] = models.Metrics{ID: "MSpanSys", MType: models.Gauge, Value: &mSpanSys}

	mallocs := float64(memStats.Mallocs)
	c.metrics["Mallocs"] = models.Metrics{ID: "Mallocs", MType: models.Gauge, Value: &mallocs}

	nextGC := float64(memStats.NextGC)
	c.metrics["NextGC"] = models.Metrics{ID: "NextGC", MType: models.Gauge, Value: &nextGC}

	numForcedGC := float64(memStats.NumForcedGC)
	c.metrics["NumForcedGC"] = models.Metrics{ID: "NumForcedGC", MType: models.Gauge, Value: &numForcedGC}

	numGC := float64(memStats.NumGC)
	c.metrics["NumGC"] = models.Metrics{ID: "NumGC", MType: models.Gauge, Value: &numGC}

	otherSys := float64(memStats.OtherSys)
	c.metrics["OtherSys"] = models.Metrics{ID: "OtherSys", MType: models.Gauge, Value: &otherSys}

	pauseTotalNs := float64(memStats.PauseTotalNs)
	c.metrics["PauseTotalNs"] = models.Metrics{ID: "PauseTotalNs", MType: models.Gauge, Value: &pauseTotalNs}

	stackInuse := float64(memStats.StackInuse)
	c.metrics["StackInuse"] = models.Metrics{ID: "StackInuse", MType: models.Gauge, Value: &stackInuse}

	stackSys := float64(memStats.StackSys)
	c.metrics["StackSys"] = models.Metrics{ID: "StackSys", MType: models.Gauge, Value: &stackSys}

	sys := float64(memStats.Sys)
	c.metrics["Sys"] = models.Metrics{ID: "Sys", MType: models.Gauge, Value: &sys}

	totalAlloc := float64(memStats.TotalAlloc)
	c.metrics["TotalAlloc"] = models.Metrics{ID: "TotalAlloc", MType: models.Gauge, Value: &totalAlloc}

	randomValue := rand.Float64()
	c.metrics["RandomValue"] = models.Metrics{ID: "RandomValue", MType: models.Gauge, Value: &randomValue}

	pollCount := c.pollCount
	c.metrics["PollCount"] = models.Metrics{ID: "PollCount", MType: models.Counter, Delta: &pollCount}

	c.pollCount++
	return c.metrics
}

func (c *Collector) SendMetricByHTTP(metric models.Metrics) {
	url := fmt.Sprintf("http://%s/update/", c.cfg.ServerAddress)
	body, err := json.Marshal(metric)
	if err != nil {
		c.logger.Errorln(err)
		return
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err = gz.Write(body)
	if err != nil {
		c.logger.Errorln(err)
		return
	}

	err = gz.Close()
	if err != nil {
		c.logger.Errorln(err)
		return
	}

	client := http.Client{}

	for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
		req, err := http.NewRequest("POST", url, &buf)
		if err != nil {
			c.logger.Errorln(err)
			return
		}

		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err == nil && res.StatusCode == http.StatusOK {
			defer res.Body.Close()
			return
		}

		if err != nil {
			c.logger.Errorln("Error sending metric:", err)
		} else if res.StatusCode != http.StatusOK {
			c.logger.Errorln("Received non-OK response:", res.Status)
			res.Body.Close()
		}

		time.Sleep(delay)
	}

	c.logger.Errorln("Failed to send metric after retries")
}

func (s *Collector) Run() {
	var metrics map[string]models.Metrics

	go func() {
		for range time.Tick(time.Duration(s.cfg.PollInterval) * time.Second) {
			metrics = s.GetMetric()
		}
	}()

	for range time.Tick(time.Duration(s.cfg.ReportInterval) * time.Second) {
		for _, m := range metrics {
			s.SendMetricByHTTP(m)
			metrics = nil
		}
	}
}
