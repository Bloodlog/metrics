package router

import (
	"github.com/Bloodlog/metrics/internal/server/handlers"
	"github.com/Bloodlog/metrics/internal/server/storage"
	"log"
	"net/http"
)

func Run() {
	memStorage := storage.NewMemStorage()

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/gauge/`, handlers.GaugeHandler(memStorage))
	mux.HandleFunc(`/update/counter/`, handlers.CounterHandler(memStorage))
	mux.HandleFunc(`/update/`, DefaultHandler)
	log.Fatal(http.ListenAndServe(`:8080`, mux))
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
