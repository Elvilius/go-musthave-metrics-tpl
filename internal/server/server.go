package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	middleware "github.com/Elvilius/go-musthave-metrics-tpl/internal/midleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	handler *handler.Handler
	router  *chi.Mux
	cfg     *config.ServerConfig
	logger  *zap.SugaredLogger
}

func New(cfg *config.ServerConfig, handler *handler.Handler, logger *zap.SugaredLogger) *Server {
	router := chi.NewRouter()

	server := &Server{handler: handler, router: router, cfg: cfg, logger: logger}

	router.Use(middleware.Logging(*logger))
	router.Use(middleware.Gzip)

	router.Get("/", server.handler.All)
	router.Post("/update/{type}/{id}/{value}", server.handler.Update)
	router.Post("/update/", server.handler.UpdateJSON)
	router.Get("/value/{type}/{id}", server.handler.Value)
	router.Post("/value/", server.handler.ValueJSON)

	return server
}

func (s *Server) Run(storage handler.Storager) {
	err := storage.LoadFromFile()
	if err != nil {
		s.logger.Errorln(err)
	}

	go func() {
		fmt.Println("Starting server...")
		err = http.ListenAndServe(s.cfg.Address, s.router)
		if err != nil {
			s.logger.Errorln(err)
			os.Exit(1)
		}
	}()

	ticker := time.NewTicker(time.Duration(s.cfg.StoreInterval) * time.Second)
	defer ticker.Stop()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ticker.C:
				err := storage.SaveToFile()
				if err != nil {
					s.logger.Errorln(err)
				}
			case <-done:
				err := storage.SaveToFile()
				if err != nil {
					s.logger.Errorln(err)
				}
				close(done)
				return
			}
		}
	}()

	<-done
}
