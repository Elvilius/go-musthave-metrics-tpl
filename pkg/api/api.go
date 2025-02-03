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

type API struct {
	url        string
	client     *http.Client
	logger     *zap.SugaredLogger
	gzipBuffer bytes.Buffer
	gzipWriter *gzip.Writer
	gzipMutex  sync.Mutex
}

func New(url string, logger *zap.SugaredLogger) *API {
	return &API{
		url:        url,
		client:     &http.Client{},
		logger:     logger,
		gzipWriter: gzip.NewWriter(&bytes.Buffer{}),
	}
}

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

		api.gzipWriter.Flush()

		compressedData := make([]byte, api.gzipBuffer.Len())
		copy(compressedData, api.gzipBuffer.Bytes())

		api.gzipWriter.Close()
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
			defer res.Body.Close()
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
