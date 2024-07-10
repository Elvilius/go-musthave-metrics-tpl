package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/gzip"
	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func Logging(logger zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			responseData := &responseData{size: 0, status: 0}

			lw := &loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}

			h.ServeHTTP(lw, r)

			duration := time.Since(start)

			logger.Infoln(
				"uri", r.RequestURI,
				"method", r.Method,
				"status", lw.responseData.status,
				"duration", duration,
				"size", lw.responseData.size,
			)
		})
	}
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.responseData.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := w

		encodings := r.Header.Get("Accept-Encoding")
		// Если принимаем gzip то сжимаем
		if strings.Contains(encodings, "gzip") {
			cw := gzip.NewCompressWriter(w)
			rw = cw

			defer cw.Close()
		}

		// Если содержимое содержит gzip
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			cr, err := gzip.NewCompressReader(r.Body)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(rw, r)
	})
}