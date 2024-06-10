package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/storage"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

type Handler struct {
	s storage.Storage
}

func NewHandler(s storage.Storage) Handler {
	return Handler{s: s}
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/update/")
	parts := make([]string, 3)
	for i, item := range strings.Split(path, "/") {
		parts[i] = item
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	metricType := parts[0]
	metricName := parts[1]
	metricValue := parts[2]

	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if metricType == "" || metricValue == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch metricType {
	case gauge:
		h.s.Gauge(metricName, value)
	case counter:
		h.s.Inc(metricName)
	default:
		{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
