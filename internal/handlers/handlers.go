package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/hashing"
	"github.com/go-chi/chi/v5"
)

type Metrics interface {
	Add(ctx context.Context, metric models.Metrics, value string) error
	GetOne(ctx context.Context, mType, id string) (*models.Metrics, error)
	GetAll(ctx context.Context) ([]models.Metrics, error)
	Update(ctx context.Context, metric []models.Metrics) error
}

type Handler struct {
	cfg     *config.ServerConfig
	metrics Metrics
}

func NewHandler(cfg *config.ServerConfig, metrics Metrics) *Handler {
	return &Handler{metrics: metrics, cfg: cfg}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	defer r.Body.Close()

	ctx := r.Context()

	mType := chi.URLParam(r, "type")
	id := chi.URLParam(r, "id")
	value := chi.URLParam(r, "value")

	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if mType == "" || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if mType != models.Counter && mType != models.Gauge {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metric := models.Metrics{
		ID:    id,
		MType: mType,
	}

	err := h.metrics.Add(ctx, metric, value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Value(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	defer r.Body.Close()

	ctx := r.Context()
	mType := chi.URLParam(r, "type")
	id := chi.URLParam(r, "id")

	metric, err := h.metrics.GetOne(ctx, mType, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if metric == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := metric.MarshalValue()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.addHash(w, bytes)

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	ctx := r.Context()
	var metric models.Metrics
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.metrics.Add(ctx, metric, "")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to update metric", http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(metric)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	h.addHash(w, res)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (h *Handler) ValueJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	metric := models.Metrics{}

	err := json.NewDecoder(r.Body).Decode(&metric)

	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	m, err := h.metrics.GetOne(ctx, metric.MType, metric.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.addHash(w, bytes)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	ctx := r.Context()

	m, err := h.metrics.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

func (h *Handler) UpdatesJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	ctx := r.Context()
	requestMetrics := []models.Metrics{}
	err := json.NewDecoder(r.Body).Decode(&requestMetrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errUpdate := h.metrics.Update(ctx, requestMetrics)

	if errUpdate != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) addHash(w http.ResponseWriter, data []byte) {
	if h.cfg.Key == "" {
		return
	}

	hash := hashing.GenerateHash(h.cfg.Key, data)
	w.Header().Set("HashSHA256", hash)
}
