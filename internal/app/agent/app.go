package agent

import (
	"context"
	"fmt"
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
	headers := make(map[string]string)
	defer wg.Done()
	for metric := range jobs {
		app.logger.Infof("Worker %d processing metric", id)
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

func (app *AppAgent) Run(ctx context.Context) {
	collectTicker := time.NewTicker(time.Duration(app.cfg.PollInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(app.cfg.ReportInterval) * time.Second)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	jobs := make(chan models.Metrics)
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(ctx)

	cleanApp := func() {
		cancel()
		wg.Wait()
		sendTicker.Stop()
		collectTicker.Stop()
		close(jobs)
		fmt.Println("Resources cleaned up")
	}

	defer cleanApp()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled")
			return
		case <-stop:
			fmt.Println("Received termination signal")
			os.Exit(1)
			return
		case <-collectTicker.C:
			go func() {
				select {
				case <-ctx.Done():
					return
				default:
					app.collector.CollectMetric()
					metrics := app.collector.GetMetrics()
					for _, m := range metrics {
						select {
						case jobs <- m:
						case <-ctx.Done():
							return
						}
					}
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
