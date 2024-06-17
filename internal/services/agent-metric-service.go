package services

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/domain"
)

const (
	Gauge   domain.MetricType = "gauge"
	Counter domain.MetricType = "counter"
)

type AgentMetricService struct {
	cfg       config.AgentConfig
	metrics   map[string]domain.Metric
	pollCount float64
}

func NewAgentMetricService(cfg config.AgentConfig) *AgentMetricService {
	return &AgentMetricService{cfg: cfg, metrics: make(map[string]domain.Metric), pollCount: 0}
}

func (s *AgentMetricService) GetMetric() map[string]domain.Metric {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	s.metrics["Alloc"] = domain.Metric{Name: "Alloc", Type: Gauge, Value: (float64(memStats.Alloc))}
	s.metrics["BuckHashSys"] = domain.Metric{Name: "BuckHashSys", Type: Gauge, Value: (float64(memStats.BuckHashSys))}
	s.metrics["Frees"] = domain.Metric{Name: "Frees", Type: Gauge, Value: (float64(memStats.Frees))}
	s.metrics["GCCPUFraction"] = domain.Metric{Name: "GCCPUFraction", Type: Gauge, Value: (memStats.GCCPUFraction)}
	s.metrics["GCSys"] = domain.Metric{Name: "GCSys", Type: Gauge, Value: (float64(memStats.GCSys))}
	s.metrics["HeapAlloc"] = domain.Metric{Name: "HeapAlloc", Type: Gauge, Value: (float64(memStats.HeapAlloc))}
	s.metrics["HeapIdle"] = domain.Metric{Name: "HeapIdle", Type: Gauge, Value: (float64(memStats.HeapIdle))}
	s.metrics["HeapInuse"] = domain.Metric{Name: "HeapInuse", Type: Gauge, Value: (float64(memStats.HeapInuse))}
	s.metrics["HeapObjects"] = domain.Metric{Name: "HeapObjects", Type: Gauge, Value: (float64(memStats.HeapObjects))}
	s.metrics["HeapReleased"] = domain.Metric{Name: "HeapReleased", Type: Gauge, Value: (float64(memStats.HeapReleased))}
	s.metrics["HeapSys"] = domain.Metric{Name: "HeapSys", Type: Gauge, Value: (float64(memStats.HeapSys))}
	s.metrics["LastGC"] = domain.Metric{Name: "LastGC", Type: Gauge, Value: (float64(memStats.LastGC))}
	s.metrics["Lookups"] = domain.Metric{Name: "Lookups", Type: Gauge, Value: (float64(memStats.Lookups))}
	s.metrics["MCacheInuse"] = domain.Metric{Name: "MCacheInuse", Type: Gauge, Value: (float64(memStats.MCacheInuse))}
	s.metrics["MCacheSys"] = domain.Metric{Name: "MCacheSys", Type: Gauge, Value: (float64(memStats.MCacheSys))}
	s.metrics["MSpanInuse"] = domain.Metric{Name: "MSpanInuse", Type: Gauge, Value: (float64(memStats.MSpanInuse))}
	s.metrics["MSpanSys"] = domain.Metric{Name: "MSpanSys", Type: Gauge, Value: (float64(memStats.MSpanSys))}
	s.metrics["Mallocs"] = domain.Metric{Name: "Mallocs", Type: Gauge, Value: (float64(memStats.Mallocs))}
	s.metrics["NextGC"] = domain.Metric{Name: "NextGC", Type: Gauge, Value: (float64(memStats.NextGC))}
	s.metrics["NumForcedGC"] = domain.Metric{Name: "NumForcedGC", Type: Gauge, Value: (float64(memStats.NumForcedGC))}
	s.metrics["NumGC"] = domain.Metric{Name: "NumGC", Type: Gauge, Value: (float64(memStats.NumGC))}
	s.metrics["OtherSys"] = domain.Metric{Name: "OtherSys", Type: Gauge, Value: (float64(memStats.OtherSys))}
	s.metrics["PauseTotalNs"] = domain.Metric{Name: "PauseTotalNs", Type: Gauge, Value: (float64(memStats.PauseTotalNs))}
	s.metrics["StackInuse"] = domain.Metric{Name: "StackInuse", Type: Gauge, Value: (float64(memStats.StackInuse))}
	s.metrics["StackSys"] = domain.Metric{Name: "StackSys", Type: Gauge, Value: (float64(memStats.StackSys))}
	s.metrics["Sys"] = domain.Metric{Name: "Sys", Type: Gauge, Value: (float64(memStats.Sys))}
	s.metrics["TotalAlloc"] = domain.Metric{Name: "TotalAlloc", Type: Gauge, Value: (float64(memStats.TotalAlloc))}
	s.metrics["RandomValue"] = domain.Metric{Name: "RandomValue", Type: Gauge, Value: (rand.Float64())}
	s.metrics["PollCount"] = domain.Metric{Name: "PollCount", Type: Counter, Value: s.pollCount}

	s.pollCount++
	return s.metrics
}

func (s *AgentMetricService) SendMetricByHttp(m domain.Metric) {
	resp, err := http.Post(fmt.Sprintf("http://%s/update/%s/%s/%f", s.cfg.ServerAddress, m.Type, m.Name, m.Value), "text/plain", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func (s *AgentMetricService) SendMetrics() {
	var metrics map[string]domain.Metric

	go func() {
		for range time.Tick(time.Duration(s.cfg.PollInterval) * time.Second) {
			metrics = s.GetMetric()
		}
	}()

	for range time.Tick(time.Duration(s.cfg.ReportInterval) * time.Second) {
		for _, m := range metrics {
			s.SendMetricByHttp(m)
			metrics = nil
		}
	}
}
