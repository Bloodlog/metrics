package router

import (
	"errors"
	"metrics/internal/server/config"
	"metrics/internal/server/handlers"
	"metrics/internal/server/handlers/counter"
	"metrics/internal/server/handlers/gauge"
	"metrics/internal/server/repository"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go.uber.org/zap"
)

func Run(configs *config.Config, memStorage *repository.MemStorage, logger zap.SugaredLogger) error {
	serverAddr := net.JoinHostPort(configs.NetAddress.Host, configs.NetAddress.Port)

	router := chi.NewRouter()

	register(router, memStorage, logger)

	logger.Infow(
		"Starting server",
		"addr", serverAddr,
	)
	err := http.ListenAndServe(serverAddr, router)
	if err != nil {
		logger.Info(err.Error(), "event", "start server")
		return errors.New("failed to start server")
	}
	return nil
}

func register(r *chi.Mux, memStorage *repository.MemStorage, logger zap.SugaredLogger) {
	r.Use(LoggingMiddleware(logger))

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

func LoggingMiddleware(logger zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			logger.Infoln(
				"method", r.Method,
				"uri", r.RequestURI,
				"status", ww.Status(),
				"size", ww.BytesWritten(),
				"duration", duration,
			)
		})
	}
}
