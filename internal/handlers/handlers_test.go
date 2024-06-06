package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/storage"
	"github.com/stretchr/testify/assert"
)

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
		memStorage := storage.NewMemStorage()
		h := NewHandler(memStorage)

		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h.Update(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.status, result.StatusCode)
		})
	}
}
