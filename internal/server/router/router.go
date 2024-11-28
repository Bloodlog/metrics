package router

import (
	"metrics/internal/server/handlers"
	"metrics/internal/server/repository"
	"net/http"
)

func Run(memStorage *repository.MemStorage, debug bool) error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/gauge/`, handlers.GaugeHandler(memStorage, debug))
	mux.HandleFunc(`/update/counter/`, handlers.CounterHandler(memStorage, debug))
	mux.HandleFunc(`/update/`, DefaultHandler)

	return http.ListenAndServe(`:8080`, mux)
}

func DefaultHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusBadRequest)
}
