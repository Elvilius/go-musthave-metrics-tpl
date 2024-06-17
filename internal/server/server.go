package server

import (
	"net/http"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	handler *handler.Handler
	r       *chi.Mux
}

func NewServer(cfg *config.ServerConfig, handler *handler.Handler) *Server {
	r := chi.NewRouter()
	return &Server{handler: handler, r: r}
}

func (s *Server) Run() {
	s.r.Get("/", s.handler.All)
	s.r.Post("/update/{metricType}/{metricName}/{metricValue}", s.handler.Update)
	s.r.Get("/value/{metricType}/{metricName}", s.handler.Value)

	cfg := config.GetServerConfig()

	err := http.ListenAndServe(cfg.Address, s.r)
	if err != nil {
		panic(err)
	}
}
