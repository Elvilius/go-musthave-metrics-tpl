package agent

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	collector "github.com/Elvilius/go-musthave-metrics-tpl/internal/collector"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/api"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/logger"
	"go.uber.org/zap"
)

type AppAgent struct {
	collector *collector.Collector
	logger    *zap.SugaredLogger
	cfg       *config.AgentConfig
	api       *api.API
}

func New() *AppAgent {
	logger, err := logger.New()
	if err != nil {
		panic(err)
	}
	cfg := config.NewAgent()

	collector := collector.New(cfg, logger)
	api := api.New(cfg.ServerAddress, logger)

	return &AppAgent{
		collector: collector,
		logger:    logger,
		cfg:       cfg,
		api:       api,
	}
}

func (app *AppAgent) Run(ctx context.Context) {
	var metrics map[string]models.Metrics

	go func() {
		for range time.Tick(time.Duration(app.cfg.PollInterval) * time.Second) {
			metrics = app.collector.GetMetric()
		}
	}()

	for range time.Tick(time.Duration(app.cfg.ReportInterval) * time.Second) {
		for _, m := range metrics {
			body, err := json.Marshal(m)

			if err != nil {
				app.logger.Fatal(err)
			}
			app.api.Fetch(ctx, http.MethodPost, "/update", body)
			metrics = nil
		}
	}
}
