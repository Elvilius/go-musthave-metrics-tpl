package server

import (
	"net/http"

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

	router.Get("/", server.handler.All)
	router.Post("/update/{type}/{id}/{value}", server.handler.Update)
	router.Post("/update/", server.handler.UpdateJSON)
	router.Post("/value", server.handler.ValueJSON)
	router.Post("/value/", server.handler.ValueJSON)
	router.Get("/value/{type}/{id}", server.handler.Value)



	return server
}

func (s *Server) Run() {
	err := http.ListenAndServe(s.cfg.Address, middleware.Logging(s.router, s.logger))
	if err != nil {
		panic(err)
	}
}
