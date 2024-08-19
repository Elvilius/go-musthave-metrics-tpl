package services

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

type Agent struct {
	cfg       *config.AgentConfig
	metrics   map[string]models.Metrics
	pollCount int64
	logger    *zap.SugaredLogger
}

func NewAgentMetricService(cfg *config.AgentConfig, logger *zap.SugaredLogger) *Agent {
	return &Agent{cfg: cfg, metrics: make(map[string]models.Metrics), pollCount: 0, logger: logger}
}

func (s *Agent) GetMetric() map[string]models.Metrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	alloc := float64(memStats.Alloc)
	s.metrics["Alloc"] = models.Metrics{ID: "Alloc", MType: models.Gauge, Value: &alloc}

	buckHashSys := float64(memStats.Frees)
	s.metrics["BuckHashSys"] = models.Metrics{ID: "BuckHashSys", MType: models.Gauge, Value: &buckHashSys}

	frees := float64(memStats.Frees)
	s.metrics["Frees"] = models.Metrics{ID: "Frees", MType: models.Gauge, Value: &frees}

	gCCPUFraction := memStats.GCCPUFraction
	s.metrics["GCCPUFraction"] = models.Metrics{ID: "GCCPUFraction", MType: models.Gauge, Value: &gCCPUFraction}

	gCSys := float64(memStats.GCSys)
	s.metrics["GCSys"] = models.Metrics{ID: "GCSys", MType: models.Gauge, Value: &gCSys}

	heapAlloc := float64(memStats.HeapAlloc)
	s.metrics["HeapAlloc"] = models.Metrics{ID: "HeapAlloc", MType: models.Gauge, Value: &heapAlloc}

	heapIdle := float64(memStats.HeapIdle)
	s.metrics["HeapIdle"] = models.Metrics{ID: "HeapIdle", MType: models.Gauge, Value: &heapIdle}

	heapInuse := float64(memStats.HeapInuse)
	s.metrics["HeapInuse"] = models.Metrics{ID: "HeapInuse", MType: models.Gauge, Value: &heapInuse}

	heapObjects := float64(memStats.HeapObjects)
	s.metrics["HeapObjects"] = models.Metrics{ID: "HeapObjects", MType: models.Gauge, Value: &heapObjects}

	heapReleased := float64(memStats.HeapReleased)
	s.metrics["HeapReleased"] = models.Metrics{ID: "HeapReleased", MType: models.Gauge, Value: &heapReleased}

	heapSys := float64(memStats.HeapSys)
	s.metrics["HeapSys"] = models.Metrics{ID: "HeapSys", MType: models.Gauge, Value: &heapSys}

	lastGC := float64(memStats.LastGC)
	s.metrics["LastGC"] = models.Metrics{ID: "LastGC", MType: models.Gauge, Value: &lastGC}

	lookups := float64(memStats.Lookups)
	s.metrics["Lookups"] = models.Metrics{ID: "Lookups", MType: models.Gauge, Value: &lookups}

	mCacheInuse := float64(memStats.MCacheInuse)
	s.metrics["MCacheInuse"] = models.Metrics{ID: "MCacheInuse", MType: models.Gauge, Value: &mCacheInuse}

	mCacheSys := float64(memStats.MCacheSys)
	s.metrics["MCacheSys"] = models.Metrics{ID: "MCacheSys", MType: models.Gauge, Value: &mCacheSys}

	mSpanInuse := float64(memStats.MSpanInuse)
	s.metrics["MSpanInuse"] = models.Metrics{ID: "MSpanInuse", MType: models.Gauge, Value: &mSpanInuse}

	mSpanSys := float64(memStats.MSpanSys)
	s.metrics["MSpanSys"] = models.Metrics{ID: "MSpanSys", MType: models.Gauge, Value: &mSpanSys}

	mallocs := float64(memStats.Mallocs)
	s.metrics["Mallocs"] = models.Metrics{ID: "Mallocs", MType: models.Gauge, Value: &mallocs}

	nextGC := float64(memStats.NextGC)
	s.metrics["NextGC"] = models.Metrics{ID: "NextGC", MType: models.Gauge, Value: &nextGC}

	numForcedGC := float64(memStats.NumForcedGC)
	s.metrics["NumForcedGC"] = models.Metrics{ID: "NumForcedGC", MType: models.Gauge, Value: &numForcedGC}

	numGC := float64(memStats.NumGC)
	s.metrics["NumGC"] = models.Metrics{ID: "NumGC", MType: models.Gauge, Value: &numGC}

	otherSys := float64(memStats.OtherSys)
	s.metrics["OtherSys"] = models.Metrics{ID: "OtherSys", MType: models.Gauge, Value: &otherSys}

	pauseTotalNs := float64(memStats.PauseTotalNs)
	s.metrics["PauseTotalNs"] = models.Metrics{ID: "PauseTotalNs", MType: models.Gauge, Value: &pauseTotalNs}

	stackInuse := float64(memStats.StackInuse)
	s.metrics["StackInuse"] = models.Metrics{ID: "StackInuse", MType: models.Gauge, Value: &stackInuse}

	stackSys := float64(memStats.StackSys)
	s.metrics["StackSys"] = models.Metrics{ID: "StackSys", MType: models.Gauge, Value: &stackSys}

	sys := float64(memStats.Sys)
	s.metrics["Sys"] = models.Metrics{ID: "Sys", MType: models.Gauge, Value: &sys}

	totalAlloc := float64(memStats.TotalAlloc)
	s.metrics["TotalAlloc"] = models.Metrics{ID: "TotalAlloc", MType: models.Gauge, Value: &totalAlloc}

	randomValue := rand.Float64()
	s.metrics["RandomValue"] = models.Metrics{ID: "RandomValue", MType: models.Gauge, Value: &randomValue}

	pollCount := s.pollCount
	s.metrics["PollCount"] = models.Metrics{ID: "PollCount", MType: models.Counter, Delta: &pollCount}

	s.pollCount++
	return s.metrics
}

func (s *Agent) SendMetricByHTTP(metric models.Metrics) {
	url := fmt.Sprintf("http://%s/update/", s.cfg.ServerAddress)
	body, err := json.Marshal(metric)
	if err != nil {
		s.logger.Errorln(err)
		return
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err = gz.Write(body)
	if err != nil {
		s.logger.Errorln(err)
		return
	}

	err = gz.Close()
	if err != nil {
		s.logger.Errorln(err)
		return
	}

	client := http.Client{}

	for _, delay := range []time.Duration{time.Second, 2*time.Second, 3*time.Second} {
		req, err := http.NewRequest("POST", url, &buf)
		if err != nil {
			s.logger.Errorln(err)
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
			s.logger.Errorln("Error sending metric:", err)
		} else if res.StatusCode != http.StatusOK {
			s.logger.Errorln("Received non-OK response:", res.Status)
			res.Body.Close()
		}

		time.Sleep(delay)
	}

	s.logger.Errorln("Failed to send metric after retries")
}

func (s *Agent) Run() {
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
