package router

import (
	"errors"
	"log"
	"metrics/internal/server/config"
	"metrics/internal/server/handlers"
	"metrics/internal/server/handlers/counter"
	"metrics/internal/server/handlers/gauge"
	"metrics/internal/server/repository"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	ErrFailedStartServer = errors.New("failed to start server")
)

func Run(configs *config.Config, memStorage *repository.MemStorage) error {
	serverAddr := net.JoinHostPort(configs.NetAddress.Host, configs.NetAddress.Port)

	router := chi.NewRouter()

	register(router, memStorage, configs.Debug)

	log.Printf("Запуск сервера на адресе: %v", serverAddr)
	err := http.ListenAndServe(serverAddr, router)
	if err != nil {
		log.Printf("failed to start server on %s: %v", serverAddr, err)
		return ErrFailedStartServer
	}
	return nil
}

func register(r *chi.Mux, memStorage *repository.MemStorage, debug bool) {
	if debug {
		r.Use(middleware.Logger)
	}

	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{metricName}/{metricValue}", gauge.UpdateGaugeHandler(memStorage))
		r.Post("/counter/{counterName}/{counterValue}", counter.UpdateCounterHandler(memStorage))
		r.Post("/{metricType}/{counterName}/{counterValue}", validateMetricType)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/gauge/{metricName}", gauge.GetGaugeHandler(memStorage))
		r.Get("/counter/{metricName}", counter.GetCounterHandler(memStorage))
		r.Get("/{metricType}/{counterName}", validateMetricType)
	})

	r.Get("/", handlers.ListHandler(memStorage))
}

func validateMetricType(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
