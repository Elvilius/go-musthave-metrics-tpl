package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	storage Storager
}

type Storager interface {
	Save(ctx context.Context, metric models.Metrics) error
	Get(ctx context.Context, mType, id string) (models.Metrics, bool, error)
	GetAll(ctx context.Context) ([]models.Metrics, error)
	Updates(ctx context.Context, metrics []models.Metrics) error
}

func NewHandler(storage Storager) *Handler {
	return &Handler{storage: storage}
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
	mType := chi.URLParam(r, "type")
	id := chi.URLParam(r, "id")
	var err error

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	m, ok, err := h.storage.Get(r.Context(), mType, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	requestMetric := models.Metrics{}
	err := json.NewDecoder(r.Body).Decode(&requestMetric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")

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
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdatesJSON(w http.ResponseWriter, r *http.Request) {
	fmt.Println(123123123123)
	requestMetrics := []models.Metrics{}
	err := json.NewDecoder(r.Body).Decode(&requestMetrics)
	fmt.Println(err)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")

	for _, re := range requestMetrics {
		fmt.Println(re)
	}
	// var responseMetric models.Metrics

	// err = h.storage.Save(r.Context(), requestMetric)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// metric, ok, err := h.storage.Get(r.Context(), requestMetric.MType, requestMetric.ID)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// if !ok {
	// 	responseMetric = requestMetric
	// } else {
	// 	responseMetric = metric
	// }

	// res, err := json.Marshal(responseMetric)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// w.WriteHeader(http.StatusOK)

	// _, err = w.Write(res)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
}
