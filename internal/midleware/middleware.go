package middleware

import (
	"net/http"
	"time"

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

func Logging(h http.Handler, logger *zap.SugaredLogger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
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
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}

	return http.HandlerFunc(fn)
}
