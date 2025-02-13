package agent

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	collector "github.com/Elvilius/go-musthave-metrics-tpl/internal/collector"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/api"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/hashing"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/logger"
	"go.uber.org/zap"
)

type AppAgent struct {
	collector *collector.Collector
	logger    *zap.SugaredLogger
	cfg       *config.AgentConfig
	api       *api.API
	sync.WaitGroup
}

func New() (*AppAgent, error) {
	logger, err := logger.New()
	if err != nil {
		return nil, err
	}
	cfg := config.NewAgent()

	collector := collector.New(cfg, logger)
	api := api.New(cfg.ServerAddress, logger)

	return &AppAgent{
		collector: collector,
		logger:    logger,
		cfg:       cfg,
		api:       api,
	}, nil
}

func (app *AppAgent) Run(ctx context.Context) {
	collectTicker := time.NewTicker(time.Duration(app.cfg.PollInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(app.cfg.ReportInterval) * time.Second)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	metricsCh := make(chan []*models.Metrics)
	sendMetricsCh := make(chan []*models.Metrics)

	ctx, cancel := context.WithCancel(ctx)
	app.RegisterWorker(ctx, sendMetricsCh)

	defer func() {
		cancel()
		sendTicker.Stop()
		collectTicker.Stop()
		app.Wait()
		close(metricsCh)
		close(sendMetricsCh)

	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-stop:
			return
		case <-collectTicker.C:
			app.Add(1)
			go func() {
				defer app.Done()
				app.collector.CollectMetric()
				metrics := app.collector.GetMetrics()
				select {
				case metricsCh <- metrics:
				case <-ctx.Done():
					return
				}
			}()
		case <-sendTicker.C:
			select {
			case metrics := <-metricsCh:
				select {
				case sendMetricsCh <- metrics:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

func (app *AppAgent) Worker(ctx context.Context, id int, jobs <-chan []*models.Metrics) {
	for metrics := range jobs {
		for _, metric := range metrics {
			headers := make(map[string]string)
			body, err := metric.MarshalMetric()
			if err != nil {
				app.logger.Error(err)
				continue
			}
			if app.cfg.Key != "" {
				headers["HashSHA256"] = hashing.GenerateHash(app.cfg.Key, body)
			}
			app.api.Fetch(ctx, http.MethodPost, "/update", body, headers)
		}
	}
}

func (app *AppAgent) RegisterWorker(ctx context.Context, jobs <-chan []*models.Metrics) {
	for i := 1; i <= app.cfg.RateLimit; i++ {
		localID := i
		go func(id int) {
			defer app.Done()
			app.Worker(ctx, id, jobs)
		}(localID)
	}
}
