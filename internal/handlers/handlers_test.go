package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type TestStorage struct {
	metrics map[string]models.Metrics
}

func (r *TestStorage) Save(metric models.Metrics) error {
	mType, ID, value, delta := metric.MType, metric.ID, metric.Value, metric.Delta

	existMetric, ok := r.Get(mType, ID)

	if mType == models.Gauge {
		r.metrics[ID] = models.Metrics{ID: ID, MType: mType, Value: value}
		return nil
	}
	if mType == models.Counter {
		if !ok {
			r.metrics[ID] = models.Metrics{ID: ID, MType: mType, Delta: delta}
			return nil
		} else {
			delta := *existMetric.Delta + *delta
			existMetric.Delta = &delta
			r.metrics[ID] = existMetric
		}
	}
	return nil
}

func (r *TestStorage) Get(mType string, ID string) (models.Metrics, bool) {
	m, ok := r.metrics[ID]
	if !ok {
		return models.Metrics{}, false
	}
	if m.MType != mType {
		return models.Metrics{}, false
	}

	return m, true
}

func (r *TestStorage) GetAll() []models.Metrics {
	all := make([]models.Metrics, 0, len(r.metrics))
	for _, m := range r.metrics {
		all = append(all, m)
	}
	return all
}

func TestHandler_Update(t *testing.T) {
	type want struct {
		status int
	}

	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name:    "positive test #1",
			request: "/update/gauge/cpu/7513",
			want: want{
				status: 200,
			},
		},
		{
			name:    "positive test #2",
			request: "/update/counter/cpu/8",
			want: want{
				status: 200,
			},
		},
		{
			name:    "negative test #1",
			request: "/update/test123123123/cpu/8",
			want: want{
				status: 400,
			},
		},
		{
			name:    "negative test #2",
			request: "/update/",
			want: want{
				status: 404,
			},
		},
	}
	for _, tt := range tests {
		memStorage := &TestStorage{metrics: make(map[string]models.Metrics)}
		h := NewHandler(memStorage)
		router := chi.NewRouter()
		router.Post("/update/{type}/{id}/{value}", h.Update)

		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.status, result.StatusCode)
			result.Body.Close()
		})
	}
}

func TestHandler_Value(t *testing.T) {
	type want struct {
		status int
	}

	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name:    "positive test #1",
			request: "/value/gauge/Alloc",
			want: want{
				status: 200,
			},
		},
		{
			name:    "positive test #2",
			request: "/value/counter/PollCount",
			want: want{
				status: 200,
			},
		},
		{
			name:    "negative test #1",
			request: "/value/test/Alloc",
			want: want{
				status: 404,
			},
		},
		{
			name:    "negative test #2",
			request: "/value/counter/test",
			want: want{
				status: 404,
			},
		},
	}
	memStorage := &TestStorage{metrics: make(map[string]models.Metrics)}

	allocValue := 1.1
	allocMetric := models.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Value: &allocValue,
	}
	err := memStorage.Save(allocMetric)
	if err != nil {
		return
	}

	var pollCountValue int64 = 100
	pollCountMetric := models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollCountValue,
	}

	err = memStorage.Save(pollCountMetric)
	if err != nil {
		return
	}

	h := NewHandler(memStorage)
	router := chi.NewRouter()
	router.Get("/value/{type}/{id}", h.Value)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.status, result.StatusCode)
			result.Body.Close()
		})
	}
}
