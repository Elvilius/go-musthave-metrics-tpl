//nolint

package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/metrics"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type TestStorage struct {
	metrics map[string]models.Metrics
}

func (r *TestStorage) Save(ctx context.Context, metric models.Metrics) error {
	mType, ID, value, delta := metric.MType, metric.ID, metric.Value, metric.Delta

	existMetric, ok, err := r.Get(context.TODO(), mType, ID)
	if err != nil {
		return err
	}

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

func (r *TestStorage) Get(ctx context.Context, mType string, ID string) (models.Metrics, bool, error) {
	m, ok := r.metrics[ID]
	if !ok {
		return models.Metrics{}, false, nil
	}
	if m.MType != mType {
		return models.Metrics{}, false, nil
	}

	return m, true, nil
}

func (r *TestStorage) GetAll(ctx context.Context) ([]models.Metrics, error) {
	all := make([]models.Metrics, 0, len(r.metrics))
	for _, m := range r.metrics {
		all = append(all, m)
	}
	return all, nil
}

func (r *TestStorage) Updates(ctx context.Context, metrics []models.Metrics) error {
	for _, metric := range metrics {
		if err := r.Save(ctx, metric); err != nil {
			return err
		}
	}
	return nil
}

func TestHandler_Update(t *testing.T) {
	type want struct {
		status int
	}

	tests := []struct {
		name    string
		request string
		want    want
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

		cfg := &config.ServerConfig{}
		metricsService := metrics.New(memStorage, nil)
		h := NewHandler(cfg, nil, metricsService)
		router := chi.NewRouter()
		router.Post("/update/{type}/{id}/{value}", h.Update)

		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.status, result.StatusCode)	
		})
	}
}

func TestHandler_Value(t *testing.T) {
	type want struct {
		status int
	}

	tests := []struct {
		name    string
		request string
		want    want
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
	metricService := metrics.New(memStorage, nil)

	allocValue := 1.1
	allocMetric := models.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Value: &allocValue,
	}
	err := memStorage.Save(context.TODO(), allocMetric)
	if err != nil {
		return
	}

	var pollCountValue int64 = 100
	pollCountMetric := models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollCountValue,
	}

	err = memStorage.Save(context.TODO(), pollCountMetric)
	if err != nil {
		return
	}

	cfg := &config.ServerConfig{}
	h := NewHandler(cfg, nil, metricService)
	router := chi.NewRouter()
	router.Get("/value/{type}/{id}", h.Value)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.status, result.StatusCode)
			err := result.Body.Close()
			fmt.Println(err)
		})
	}
}

func TestHandler_GetAll(t *testing.T) {
	type want struct {
		status int
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name: "positive test #1",
			want: want{
				status: 200,
			},
		},
	}
	memStorage := &TestStorage{metrics: make(map[string]models.Metrics)}
	metricService := metrics.New(memStorage, nil)

	allocValue := 1.1
	var pollCountValue int64 = 100

	metrics := []models.Metrics{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: &allocValue,
		},
		{
			ID:    "PollCount",
			MType: "counter",
			Delta: &pollCountValue,
		},
	}

	for _, m := range metrics {
		fmt.Println(memStorage.Save(context.TODO(), m))

	}

	cfg := &config.ServerConfig{}
	h := NewHandler(cfg, nil, metricService)
	router := chi.NewRouter()
	router.Get("/", h.All)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.status, result.StatusCode)
			fmt.Println(result.Body.Close())
		})
	}
}

// ExampleHandler_Update demonstrates how to use the Update handler.
func ExampleHandler_Update() {
	memStorage := &TestStorage{metrics: make(map[string]models.Metrics)}
	cfg := &config.ServerConfig{}
	metricsService := metrics.New(memStorage, nil)
	h := NewHandler(cfg, nil, metricsService)

	router := chi.NewRouter()
	router.Post("/update/{type}/{id}/{value}", h.Update)

	request := httptest.NewRequest(http.MethodPost, "/update/gauge/cpu/7513", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)
	result := w.Result()
	defer func() {
		fmt.Println(result.Body.Close())
	}()

	fmt.Println("Status:", result.StatusCode)
	// Output:
	// Status: 200
}
