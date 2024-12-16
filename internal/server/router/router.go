package router

import (
	"errors"
	"metrics/internal/server/config"
	"metrics/internal/server/handlers/api"
	"metrics/internal/server/handlers/web"
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
		r.Post("/", api.UpdateHandler(memStorage, logger))
		r.Post("/{metricType}/{metricName}/{metricValue}", web.UpdateHandler(memStorage, logger))
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", api.GetHandler(memStorage, logger))
		r.Get("/{metricType}/{metricName}", web.GetHandler(memStorage, logger))
	})
	r.Get("/", web.ListHandler(memStorage, logger))
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
