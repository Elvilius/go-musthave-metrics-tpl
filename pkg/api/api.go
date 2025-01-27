package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type API struct {
	url    string
	client *http.Client
	logger *zap.SugaredLogger
}

func New(url string, logger *zap.SugaredLogger) *API {
	return &API{
		url:    url,
		client: &http.Client{},
		logger: logger,
	}
}

func (api *API) Fetch(ctx context.Context, method string, endpoint string, body []byte, headers map[string]string) {
	url := fmt.Sprintf("http://%s%s/", api.url, endpoint)
	client := http.Client{}
	for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)

		_, err := gz.Write(body)
		if err != nil {
			api.logger.Errorln("Error writing to gzip writer:", err)
			return
		}

		err = gz.Close()
		if err != nil {
			api.logger.Errorln("Error closing gzip writer:", err)
			return
		}

		req, err := http.NewRequest("POST", url, &buf)
		if err != nil {
			api.logger.Errorln("Error creating request:", err)
			return
		}

		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")
		for key, value := range headers {
			req.Header.Set(key, value)
		}
		res, err := client.Do(req)
		if err == nil && res.StatusCode == http.StatusOK {
			defer res.Body.Close()
			return
		}

		if err != nil {
			api.logger.Errorln("Error sending metric:", err)
		} else if res.StatusCode != http.StatusOK {
			api.logger.Errorln("Received non-OK response:", res.Status)
			res.Body.Close()
		}

		time.Sleep(delay)
	}

	api.logger.Errorln("Failed to send metric after retries")
}
