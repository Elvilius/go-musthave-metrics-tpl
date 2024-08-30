package services

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/hashing"
	"go.uber.org/zap"
)

type Agent struct {
	cfg       *config.AgentConfig
	metrics   []models.IMetric
	pollCount int64
	logger    *zap.SugaredLogger
}

func NewAgentMetricService(cfg *config.AgentConfig, logger *zap.SugaredLogger) *Agent {
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

	return &Agent{cfg: cfg, metrics: metrics, pollCount: 0, logger: logger}
}

func (s *Agent) CollectMetric() {
	for _, m := range s.metrics {
		m.Update()
	}
}

func (s *Agent) GetMetrics() []models.Metrics {
	metrics := make([]models.Metrics, len(s.metrics))
	for i, metric := range s.metrics {
		metrics[i] = metric.Get()
	}
	return metrics
}

func (s *Agent) SendMetricByHTTP(metric models.Metrics) {
	url := fmt.Sprintf("http://%s/update/", s.cfg.ServerAddress)
	body, err := json.Marshal(metric)
	if err != nil {
		s.logger.Errorln(err)
		return
	}
	client := http.Client{}
	for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)

		_, err := gz.Write(body)
		if err != nil {
			s.logger.Errorln("Error writing to gzip writer:", err)
			return
		}

		err = gz.Close()
		if err != nil {
			s.logger.Errorln("Error closing gzip writer:", err)
			return
		}

		req, err := http.NewRequest("POST", url, &buf)
		if err != nil {
			s.logger.Errorln("Error creating request:", err)
			return
		}

		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")
		if s.cfg.Key != "" {
			dataHash := hashing.GenerateHash(s.cfg.Key, body)
			req.Header.Set("HashSHA256", dataHash)
		}

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

func (s *Agent) SendMetricsByHTTP() {
	metrics := s.GetMetrics()
	for _, metric := range metrics {
		s.SendMetricByHTTP(metric)
	}
}

func (s *Agent) Worker(id int, jobs <-chan models.Metrics) {

	for metric := range jobs {
		s.logger.Infoln(fmt.Sprintf("Worker %d processing metric", id))
		s.SendMetricByHTTP(metric)
	}
}

func (s *Agent) Run(ctx context.Context) {
	collectTicker := time.NewTicker(time.Duration(s.cfg.PollInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(s.cfg.ReportInterval) * time.Second)

	for {
		select {
		case <-ctx.Done():
			collectTicker.Stop()
			sendTicker.Stop()
			return
		case <-collectTicker.C:
			go s.CollectMetric()
		case <-sendTicker.C:
			metrics := s.GetMetrics()
			jobs := make(chan models.Metrics, len(metrics))

			for i := 1; i <= s.cfg.RateLimit; i++ {
				go s.Worker(i, jobs)
			}
			for _, m := range metrics {
				jobs <- m
			}

			close(jobs)
		}
	}
}
