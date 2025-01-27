package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	middleware "github.com/Elvilius/go-musthave-metrics-tpl/pkg/midleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	handler *handler.Handler
	router  *chi.Mux
	cfg     *config.ServerConfig
	logger  *zap.SugaredLogger
}

func New(cfg *config.ServerConfig, handler *handler.Handler, logger *zap.SugaredLogger, db *sql.DB) *Server {
	router := chi.NewRouter()

	server := &Server{handler: handler, router: router, cfg: cfg, logger: logger}

	router.Use(middleware.Logging(logger))
	router.Use(middleware.Gzip)

	router.Get("/", server.handler.All)
	router.Post("/update/{type}/{id}/{value}", server.handler.Update)
	router.Post("/update/", server.handler.UpdateJSON)
	router.Get("/value/{type}/{id}", server.handler.Value)
	router.Post("/value/", server.handler.ValueJSON)
	router.Post("/updates/", middleware.VerifyHash(cfg, *logger, http.HandlerFunc(server.handler.UpdatesJSON)))

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		if db == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err := db.Ping()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return server
}

func (s *Server) Run(ctx context.Context) {
	go func() {
		fmt.Println("Starting server...")
		err := http.ListenAndServe(s.cfg.Address, s.router)
		if err != nil {
			s.logger.Errorln(err)
			os.Exit(1)
		}
	}()
	<-ctx.Done()
}
