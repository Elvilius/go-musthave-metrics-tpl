package server

import (
	"net/http"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	handler *handler.Handler
	router  *chi.Mux
	cfg     *config.ServerConfig
}

func New(cfg *config.ServerConfig, handler *handler.Handler) *Server {
	router := chi.NewRouter()

	server := &Server{handler: handler, router: router, cfg: cfg}

	router.Get("/", server.handler.All)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", server.handler.Update)
	router.Get("/value/{metricType}/{metricName}", server.handler.Value)

	return server
}

func (s *Server) Run() {
	err := http.ListenAndServe(s.cfg.Address, s.router)
	if err != nil {
		panic(err)
	}
}
