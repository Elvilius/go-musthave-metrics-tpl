package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/hashing"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	storage Storager
	cfg     *config.ServerConfig
}

type Storager interface {
	Save(ctx context.Context, metric models.Metrics) error
	Get(ctx context.Context, mType, id string) (models.Metrics, bool, error)
	GetAll(ctx context.Context) ([]models.Metrics, error)
	Updates(ctx context.Context, metrics []models.Metrics) error
}

func NewHandler(storage Storager, cfg *config.ServerConfig) *Handler {
	return &Handler{storage: storage, cfg: cfg}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	metric := models.Metrics{
		ID:    id,
		MType: mType,
	}

	if mType == models.Counter {
		parseInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric.Delta = &parseInt
	} else {
		parseFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric.Value = &parseFloat
	}

	err := h.storage.Save(r.Context(), metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Value(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	mType := chi.URLParam(r, "type")
	id := chi.URLParam(r, "id")
	var err error

	m, ok, err := h.storage.Get(r.Context(), mType, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var bytes []byte
	if m.MType == models.Counter {
		bytes, err = json.Marshal(m.Delta)
	} else {
		bytes, err = json.Marshal(m.Value)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.setHeaderHash(w, bytes)
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bytesReq, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	h.verifyRequestHash(w, r.Header.Get("HashSHA256"), bytesReq)

	bodyReader := bytes.NewReader(bytesReq)
	requestMetric := models.Metrics{}
	err = json.NewDecoder(bodyReader).Decode(&requestMetric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var responseMetric models.Metrics

	err = h.storage.Save(r.Context(), requestMetric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metric, ok, err := h.storage.Get(r.Context(), requestMetric.MType, requestMetric.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		responseMetric = requestMetric
	} else {
		responseMetric = metric
	}

	res, err := json.Marshal(responseMetric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.setHeaderHash(w, res)
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) ValueJSON(w http.ResponseWriter, r *http.Request) {
	metric := models.Metrics{}

	err := json.NewDecoder(r.Body).Decode(&metric)

	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	m, ok, err := h.storage.Get(r.Context(), metric.MType, metric.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.setHeaderHash(w, bytes)
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	m, err := h.storage.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	h.setHeaderHash(w, bytes)
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdatesJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bytesReq, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	h.verifyRequestHash(w, r.Header.Get("HashSHA256"), bytesReq)

	requestMetrics := []models.Metrics{}
	bodyReader := bytes.NewReader(bytesReq)
	err = json.NewDecoder(bodyReader).Decode(&requestMetrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	errUpdate := h.storage.Updates(ctx, requestMetrics)

	if errUpdate != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) verifyRequestHash(w http.ResponseWriter, hash string, data []byte) {
	if h.cfg.Key == "" {
		return
	}
	ok := hashing.VerifyHash(h.cfg.Key, data, hash)
	
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) setHeaderHash(w http.ResponseWriter, data []byte) {
	if h.cfg.Key == "" {
		return
	}

	hash := hashing.GenerateHash(h.cfg.Key, data)
	w.Header().Set("HashSHA256", hash)
}
