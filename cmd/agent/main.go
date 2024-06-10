package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

const (
	pollInterval   = 2
	reportInterval = 10
	serverAddress  = "http://localhost:8080"
)

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Metric struct {
	Name  string
	MType MetricType
	Value float64
}

var pollCount float64

func collectMetrics() map[string]Metric {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := make(map[string]Metric)
	metrics["Alloc"] = Metric{Name: "Alloc", MType: Gauge, Value: (float64(memStats.Alloc))}
	metrics["BuckHashSys"] = Metric{Name: "BuckHashSys", MType: Gauge, Value: (float64(memStats.BuckHashSys))}
	metrics["Frees"] = Metric{Name: "Frees", MType: Gauge, Value: (float64(memStats.Frees))}
	metrics["GCCPUFraction"] = Metric{Name: "GCCPUFraction", MType: Gauge, Value: (memStats.GCCPUFraction)}
	metrics["GCSys"] = Metric{Name: "GCSys", MType: Gauge, Value: (float64(memStats.GCSys))}
	metrics["HeapAlloc"] = Metric{Name: "HeapAlloc", MType: Gauge, Value: (float64(memStats.HeapAlloc))}
	metrics["HeapIdle"] = Metric{Name: "HeapIdle", MType: Gauge, Value: (float64(memStats.HeapIdle))}
	metrics["HeapInuse"] = Metric{Name: "HeapInuse", MType: Gauge, Value: (float64(memStats.HeapInuse))}
	metrics["HeapObjects"] = Metric{Name: "HeapObjects", MType: Gauge, Value: (float64(memStats.HeapObjects))}
	metrics["HeapReleased"] = Metric{Name: "HeapReleased", MType: Gauge, Value: (float64(memStats.HeapReleased))}
	metrics["HeapSys"] = Metric{Name: "HeapSys", MType: Gauge, Value: (float64(memStats.HeapSys))}
	metrics["LastGC"] = Metric{Name: "LastGC", MType: Gauge, Value: (float64(memStats.LastGC))}
	metrics["Lookups"] = Metric{Name: "Lookups", MType: Gauge, Value: (float64(memStats.Lookups))}
	metrics["MCacheInuse"] = Metric{Name: "MCacheInuse", MType: Gauge, Value: (float64(memStats.MCacheInuse))}
	metrics["MCacheSys"] = Metric{Name: "MCacheSys", MType: Gauge, Value: (float64(memStats.MCacheSys))}
	metrics["MSpanInuse"] = Metric{Name: "MSpanInuse", MType: Gauge, Value: (float64(memStats.MSpanInuse))}
	metrics["MSpanSys"] = Metric{Name: "MSpanSys", MType: Gauge, Value: (float64(memStats.MSpanSys))}
	metrics["Mallocs"] = Metric{Name: "Mallocs", MType: Gauge, Value: (float64(memStats.Mallocs))}
	metrics["NextGC"] = Metric{Name: "NextGC", MType: Gauge, Value: (float64(memStats.NextGC))}
	metrics["NumForcedGC"] = Metric{Name: "NumForcedGC", MType: Gauge, Value: (float64(memStats.NumForcedGC))}
	metrics["NumGC"] = Metric{Name: "NumGC", MType: Gauge, Value: (float64(memStats.NumGC))}
	metrics["OtherSys"] = Metric{Name: "OtherSys", MType: Gauge, Value: (float64(memStats.OtherSys))}
	metrics["PauseTotalNs"] = Metric{Name: "PauseTotalNs", MType: Gauge, Value: (float64(memStats.PauseTotalNs))}
	metrics["StackInuse"] = Metric{Name: "StackInuse", MType: Gauge, Value: (float64(memStats.StackInuse))}
	metrics["StackSys"] = Metric{Name: "StackSys", MType: Gauge, Value: (float64(memStats.StackSys))}
	metrics["Sys"] = Metric{Name: "Sys", MType: Gauge, Value: (float64(memStats.Sys))}
	metrics["TotalAlloc"] = Metric{Name: "TotalAlloc", MType: Gauge, Value: (float64(memStats.TotalAlloc))}
	metrics["RandomValue"] = Metric{Name: "RandomValue", MType: Gauge, Value: (rand.Float64())}
	metrics["PollCount"] = Metric{Name: "PollCount", MType: Counter, Value: pollCount}

	pollCount++
	return metrics
}

func sendMetric(metric Metric) {
	resp, err := http.Post(fmt.Sprintf("http://localhost:8080/update/%s/%s/%f", metric.MType, metric.Name, metric.Value), "text/plain", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func main() {
	var metrics map[string]Metric

	go func() {
		for range time.Tick(pollInterval * time.Second) {
			metrics = collectMetrics()
		}
	}()

	for range time.Tick(reportInterval * time.Second) {
		for _, m := range metrics {
			sendMetric(m)
			metrics = nil
		}
	}
}
