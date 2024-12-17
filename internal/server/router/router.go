package router

import (
	"errors"
	"metrics/internal/server/config"
	"metrics/internal/server/handlers/api"
	"metrics/internal/server/handlers/web"
	"metrics/internal/server/middleware"
	"metrics/internal/server/repository"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"
)

func Run(configs *config.Config, memStorage repository.MetricStorage, logger zap.SugaredLogger) error {
	serverAddr := net.JoinHostPort(configs.NetAddress.Host, configs.NetAddress.Port)

	router := chi.NewRouter()

	register(router, memStorage, logger)

	logger.Infow(
		"Starting server",
		"addr", serverAddr,
	)
	err := http.ListenAndServe(serverAddr, router)
	if err != nil {
		logger.Info(err.Error(), "router", "start server")
		return errors.New("failed to start server")
	}
	return nil
}

func register(r *chi.Mux, memStorage repository.MetricStorage, logger zap.SugaredLogger) {
	r.Use(middleware.LoggingMiddleware(logger), middleware.CompressionMiddleware(logger))

	r.Route("/update", func(r chi.Router) {
		r.Post("/", api.UpdateHandler(memStorage, logger))
		r.Post("/{metricType}/{metricName}/{metricValue}", web.UpdateHandler(memStorage, logger))
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", api.GetHandler(memStorage, logger))
		r.Get("/{metricType}/{metricName}", web.GetHandler(memStorage, logger))
	})
	r.Get("/", web.ListHandler(memStorage, logger))
}
