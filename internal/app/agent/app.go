package agent

import (
	"context"
	"fmt"
	"net/http"
	"sync"
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

func (app *AppAgent) Worker(ctx context.Context, id int, jobs <-chan models.Metrics, wg *sync.WaitGroup) {
	defer wg.Done()
	for metric := range jobs {
		app.logger.Infoln(fmt.Sprintf("Worker %d processing metric", id))
		body, err := metric.MarshalMetric()
		if err != nil {
			app.logger.Fatal(err)
		}
		app.api.Fetch(ctx, http.MethodPost, "/update", body)
	}
}

func (app *AppAgent) Run(ctx context.Context) {
	collectTicker := time.NewTicker(time.Duration(app.cfg.PollInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(app.cfg.ReportInterval) * time.Second)

	jobs := make(chan models.Metrics)
	wg := &sync.WaitGroup{}

	for {
		select {
		case <-ctx.Done():
			collectTicker.Stop()
			sendTicker.Stop()
			wg.Wait()
			close(jobs)
			return
		case <-collectTicker.C:
			go func() {
				app.collector.CollectMetric()
				metrics := app.collector.GetMetrics()

				for _, m := range metrics {
					jobs <- m
				}
			}()
		case <-sendTicker.C:
			for i := 1; i <= app.cfg.RateLimit; i++ {
				wg.Add(1)
				go app.Worker(ctx, i, jobs, wg)
			}
		}
	}
}
