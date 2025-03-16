package agent

import (
	"context"
	"net/http"
	"sync"
	"time"

	collector "github.com/Elvilius/go-musthave-metrics-tpl/internal/collector"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/api"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/crypto"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/hashing"
	"go.uber.org/zap"
)

type AppAgent struct {
	collector *collector.Collector
	logger    *zap.SugaredLogger
	cfg       *config.AgentConfig
	api       *api.API
	crypto    *crypto.Crypto
	sync.WaitGroup
}

func New(logger *zap.SugaredLogger) (*AppAgent, error) {
	cfg, err := config.NewAgent(logger)
	if err != nil {
		return nil, err
	}

	collector := collector.New(cfg, logger)
	api := api.New(cfg.ServerAddress, logger)
	crypto, err := crypto.New(crypto.Cfg{PublicKeyPath: cfg.CryptoKey})
	if err != nil {
		return nil, err
	}

	return &AppAgent{
		collector: collector,
		logger:    logger,
		cfg:       cfg,
		api:       api,
		crypto:    crypto,
	}, nil
}

func (app *AppAgent) Run(ctx context.Context) {
	collectTicker := time.NewTicker(time.Duration(app.cfg.PollInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(app.cfg.ReportInterval) * time.Second)
	metricsCh := make(chan []*models.Metrics)
	sendMetricsCh := make(chan []*models.Metrics)

	ctx, cancel := context.WithCancel(ctx)
	app.RegisterWorker(ctx, sendMetricsCh)

	defer func() {
		cancel()
		close(metricsCh)
		close(sendMetricsCh)
	}()

	app.Add(1)
	go func() {
		defer app.Done()
		for range collectTicker.C {
			select {
			case <-ctx.Done():
				collectTicker.Stop()
				return
			default:
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
			}
		}
	}()

	app.Add(1)
	go func() {
		defer app.Done()
		for range sendTicker.C {
			select {
			case <-ctx.Done():
				sendTicker.Stop()
				return
			default:
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
	}()

	app.Wait()
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
			if app.cfg.CryptoKey != "" {
				body, err = app.crypto.Encrypt(body)
				if err != nil {
					app.logger.Error(err)
					continue
				}
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
