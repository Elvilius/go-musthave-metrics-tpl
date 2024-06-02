package main

import (
	"net/http"

	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/repo"
)


func main() {
    mux := http.NewServeMux()

	repo := repo.NewRepo()
	handler := handler.NewHandler(repo)

	mux.HandleFunc("/update/", handler.Update)

	http.ListenAndServe(":8080", mux)
}
