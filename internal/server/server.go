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

	//router.Use(middleware.Logging(*logger))
	router.Use(middleware.Gzip)

	router.Get("/", server.handler.All)
	router.Post("/update/{type}/{id}/{value}", server.handler.Update)
	router.Post("/update/", server.handler.UpdateJSON)
	router.Get("/value/{type}/{id}", server.handler.Value)
	router.Post("/value/", server.handler.ValueJSON)

	return server
}

func (s *Server) Run() {
	err := http.ListenAndServe(s.cfg.Address, s.router)
	if err != nil {
		panic(err)
	}
}
