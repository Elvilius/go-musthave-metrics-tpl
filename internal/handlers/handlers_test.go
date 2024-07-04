package handler

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type TestStorage struct {
	metrics map[string]models.Metrics
}

func (r *TestStorage) Save(metricType string, metricName string, value any) error {
	existMetric, ok := r.Get(metricType, metricName)

	if metricType == models.Gauge {
		parsedValueFloat, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return err
		}
		r.metrics[metricName] = models.Metrics{MType: metricType, ID: metricName, Value: &parsedValueFloat}
		return nil
	}
	if metricType == models.Counter {
		parsedValue, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return err
		}

		if !ok {
			r.metrics[metricName] = models.Metrics{MType: metricType, ID: metricName, Delta: &parsedValue}
			return nil
		} else {
			delta := *existMetric.Delta + parsedValue
			existMetric.Delta = &delta
			r.metrics[metricName] = existMetric
		}
	}
	return nil
}

func (r *TestStorage) Get(metricType string, metricName string) (models.Metrics, bool) {
	m, ok := r.metrics[metricName]
	if !ok {
		return models.Metrics{}, false
	}
	if m.MType != metricType {
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
		router.Post("/update/{metricType}/{metricName}/{metricValue}", h.Update)

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

	err := memStorage.Save("gauge", "Alloc", "1.1")
	if err != nil {
		return
	}
	err = memStorage.Save("counter", "PollCount", "100")
	if err != nil {
		return
	}

	h := NewHandler(memStorage)
	router := chi.NewRouter()
	router.Get("/value/{metricType}/{metricName}", h.Value)
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
