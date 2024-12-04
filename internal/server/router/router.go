package router

import (
	"fmt"
	"metrics/internal/server/config"
	"metrics/internal/server/handlers"
	"metrics/internal/server/handlers/counter"
	"metrics/internal/server/handlers/gauge"
	"metrics/internal/server/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run(configs *config.Config, memStorage *repository.MemStorage) error {
	serverAddr := fmt.Sprintf("%s:%d", configs.NetAddress.Host, configs.NetAddress.Port)

	router := chi.NewRouter()

	register(router, memStorage, configs.Debug)

	fmt.Println("Запуск сервера на адресе:", serverAddr)
	return http.ListenAndServe(serverAddr, router)
}

func register(r *chi.Mux, memStorage *repository.MemStorage, debug bool) *chi.Mux {
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

	return r
}

func validateMetricType(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
