package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/metrics"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/storage"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/logger"
	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type AppServer struct {
	handler *handler.Handler
	router  *chi.Mux
	cfg     *config.ServerConfig
	logger  *zap.SugaredLogger
	store   *storage.Store
}

func New() *AppServer {
	logger, err := logger.New()
	if err != nil {
		logger.Fatal(err)
	}

	cfg := config.NewServer()

	db, err := sql.Open("postgres", cfg.DatabaseDsn)
	if err != nil {
		logger.Fatalw("Failed to open DB", "error", err)
	}
	mainStore := storage.New(cfg, logger, db)
	metricsService := metrics.New(mainStore.GetStorage(), logger)
	handler := handler.NewHandler(cfg, metricsService)

	router := chi.NewRouter()

	server := &AppServer{
		handler: handler,
		router:  router,
		cfg:     cfg,
		logger:  logger,
		store:   mainStore,
	}

	return server
}

func (a *AppServer) registerRoute() {
	m := middleware.New(a.cfg, a.logger)

	a.router.Use(m.Logging)
	a.router.Use(middleware.Gzip)
	//a.router.Use(m.VerifyHash)

	a.router.Get("/", a.handler.All)
	a.router.Post("/update/{type}/{id}/{value}", a.handler.Update)
	a.router.Post("/update/", a.handler.UpdateJSON)
	a.router.Get("/value/{type}/{id}", a.handler.Value)
	a.router.Post("/value/", a.handler.ValueJSON)
	a.router.Post("/updates/", a.handler.UpdatesJSON)

	a.router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		err := a.store.Ping()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (a *AppServer) Run(ctx context.Context) {
	a.registerRoute()
	a.store.Run(ctx)

	defer a.store.Close()
	go func() {
		fmt.Println("Starting server...")
		err := http.ListenAndServe(a.cfg.Address, a.router)
		if err != nil {
			a.logger.Errorln(err)
			os.Exit(1)
		}
	}()
	<-ctx.Done()
}
