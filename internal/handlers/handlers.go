package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/domain"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	s storage.Storage
}

func NewHandler(s storage.Storage) Handler {
	return Handler{s: s}
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if metricType == "" || metricValue == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if metricType != domain.Counter && metricType != domain.Gauge {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	err := h.s.Save(metricType, metricName, metricValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}

func (h Handler) Value(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	m, ok := h.s.Get(metricType, metricName)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(m.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) All(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	m := h.s.GetAll()

	bytes, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
