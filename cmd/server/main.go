package main

import (
	"flag"
	"net/http"

	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	memStorage := storage.NewMemStorage()
	handler := handler.NewHandler(memStorage)

	r.Get("/", handler.All)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", handler.Update)
	r.Get("/value/{metricType}/{metricName}", handler.Value)

	address := flag.String("a", "http://localhost:8080", "address")
	flag.Parse()

	err := http.ListenAndServe(*address, r)
	if err != nil {
		panic(err)
	}
}
