package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/repo"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

type Handler struct {
	repo *repo.Repo
}

func NewHandler(repo *repo.Repo) Handler {
	return Handler{repo: repo}
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
		h.repo.Gauge(metricName, value)
	case counter:
		h.repo.Inc(metricName)
	default:
		{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
