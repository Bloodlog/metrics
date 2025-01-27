package router

import (
	"fmt"
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

func Run(configs *config.Config, memStorage repository.MetricStorage, logger *zap.SugaredLogger) error {
	handlerLogger := logger.With("router", "router")
	serverAddr := net.JoinHostPort(configs.NetAddress.Host, configs.NetAddress.Port)

	router := chi.NewRouter()

	register(router, configs, memStorage, logger)

	handlerLogger.Infow(
		"Starting server",
		"addr", serverAddr,
	)
	err := http.ListenAndServe(serverAddr, router)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func register(r *chi.Mux, configs *config.Config, memStorage repository.MetricStorage, logger *zap.SugaredLogger) {
	apiHandler := api.NewHandler(memStorage, logger)
	webHandler := web.NewHandler(memStorage, logger)
	r.Use(middleware.LoggingMiddleware(logger), middleware.CompressionMiddleware(logger))

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", apiHandler.UpdatesHandler())
	})
	r.Route("/update", func(r chi.Router) {
		r.Post("/", apiHandler.UpdateHandler())
		r.Post("/{metricType}/{metricName}/{metricValue}", webHandler.UpdateHandler())
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", apiHandler.GetHandler())
		r.Get("/{metricType}/{metricName}", webHandler.GetHandler())
	})
	r.Get("/", webHandler.ListHandler())
	r.Get("/ping", webHandler.HealthHandler(configs.DatabaseDsn))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		handlerLogger := logger.With("router", "NotFound")
		handlerLogger.Infoln("Route not found",
			"method", r.Method,
			"uri", r.RequestURI,
		)
		w.WriteHeader(http.StatusNotFound)
	})
}
