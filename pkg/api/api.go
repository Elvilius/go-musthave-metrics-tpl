// Package api provides a client for sending compressed HTTP requests.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

// API represents a client for sending HTTP requests with gzip compression.
type API struct {
	url        string
	client     *http.Client
	logger     *zap.SugaredLogger
	gzipWriter *gzip.Writer
	gzipBuffer bytes.Buffer
	gzipMutex  sync.Mutex
}

// New creates a new API client with the specified URL and logger.
//
// The client is initialized with an HTTP client and a gzip writer.
func New(url string, logger *zap.SugaredLogger) *API {
	return &API{
		url:        url,
		client:     &http.Client{},
		logger:     logger,
		gzipWriter: gzip.NewWriter(&bytes.Buffer{}),
	}
}

// Fetch sends an HTTP request to the specified endpoint using the given method.
//
// The request body is compressed using gzip before sending. If the request
// fails, it is retried up to three times with increasing delays.
//
// Parameters:
//   - ctx: Context for request cancellation.
//   - method: HTTP method (e.g., "POST").
//   - endpoint: API endpoint (appended to base URL).
//   - body: Request payload (will be compressed).
//   - headers: Additional headers to include in the request.
//
// The function retries failed requests with exponential backoff.
// Logs an error if the request ultimately fails.
func (api *API) Fetch(ctx context.Context, method string, endpoint string, body []byte, headers map[string]string) {
	url := fmt.Sprintf("http://%s%s/", api.url, endpoint)

	for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
		api.gzipMutex.Lock()

		api.gzipBuffer.Reset()
		api.gzipWriter.Reset(&api.gzipBuffer)

		_, err := api.gzipWriter.Write(body)
		if err != nil {
			api.gzipMutex.Unlock()
			api.logger.Errorln("Error writing to gzip writer:", err)
			return
		}

		if errFlush := api.gzipWriter.Flush(); errFlush != nil {
			api.logger.Errorln("Error writing to gzip writer:", errFlush)
			return
		}

		compressedData := make([]byte, api.gzipBuffer.Len())
		copy(compressedData, api.gzipBuffer.Bytes())

		if errWriter := api.gzipWriter.Close(); errWriter != nil {
			api.logger.Errorln("Error close to gzip writer:", errWriter)
			return
		}

		api.gzipMutex.Unlock()

		if len(compressedData) == 0 {
			api.logger.Errorln("Error: compressedData is empty after gzipWriter.Close()")
			return
		}

		reqBody := bytes.NewReader(compressedData)

		req, err := http.NewRequest("POST", url, reqBody)
		if err != nil {
			api.logger.Errorln("Error creating request:", err)
			return
		}

		req.ContentLength = int64(len(compressedData))

		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		res, err := api.client.Do(req)
		if err == nil {
			defer func() {
				if errBodyClose := res.Body.Close(); errBodyClose != nil {
					api.logger.Errorln("Error close body", errBodyClose)
				}
			}()
			if res.StatusCode == http.StatusOK {
				return
			}
		}

		if err != nil {
			api.logger.Errorln("Error sending metric:", err)
		} else {
			api.logger.Errorln("Received non-OK response:", res.Status)
		}

		time.Sleep(delay)
	}

	api.logger.Errorln("Failed to send metric after retries")
}
