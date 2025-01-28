package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/gzip"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/hashing"
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

type Middleware struct {
	cfg    *config.ServerConfig
	logger *zap.SugaredLogger
}

func New(cfg *config.ServerConfig, logger *zap.SugaredLogger) *Middleware {
	return &Middleware{
		cfg:    cfg,
		logger: logger,
	}
}

func (m *Middleware) Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{size: 0, status: 0}

		lw := &loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(lw, r)

		duration := time.Since(start)

		m.logger.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", lw.responseData.status,
			"duration", duration,
			"size", lw.responseData.size,
		)
	})
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

func Gzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Del("Content-Length")

			gz := gzip.NewCompressWriter(w)
			defer gz.Close()

			gzipWriter := &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gz,
			}

			h.ServeHTTP(gzipWriter, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func (m *Middleware) VerifyHash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if m.cfg.Key != "" {
			if ok := hashing.VerifyHash(m.cfg.Key, data, r.Header.Get("HashSHA256")); !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		r.Body = io.NopCloser(bytes.NewBuffer(data))
		h.ServeHTTP(ow, r)

	})
}
