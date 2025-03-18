package router

import (
	"metrics/internal/server/config"
	"metrics/internal/server/handlers/api"
	"metrics/internal/server/handlers/web"
	"metrics/internal/server/middleware"
	"metrics/internal/server/repository"
	"metrics/internal/server/service"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "metrics/swagger"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func ConfigureServerHandler(
	memStorage repository.MetricStorage,
	cfg *config.Config,
	logger *zap.SugaredLogger,
) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.LoggingMiddleware(logger),
		middleware.DecompressionMiddleware(logger),
		middleware.CheckHashMiddleware(logger, cfg.Key),
		middleware.ResponseHashMiddleware(cfg.Key),
		middleware.ResponseCompressionMiddleware(logger),
	)

	register(router, cfg, memStorage, logger)

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		handlerLogger := logger.With("router", "NotFound")
		handlerLogger.Infoln("Route not found",
			"method", r.Method,
			"uri", r.RequestURI,
		)
		w.WriteHeader(http.StatusNotFound)
	})

	return router
}

func register(
	r *chi.Mux,
	cfg *config.Config,
	memStorage repository.MetricStorage,
	logger *zap.SugaredLogger,
) {
	metricService := service.NewMetricService(memStorage, logger)
	apiHandler := api.NewHandler(metricService, logger)
	webHandler := web.NewHandler(metricService, logger)

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
	r.Get("/ping", webHandler.HealthHandler(cfg.DatabaseDsn))
	r.Get("/swagger/*", httpSwagger.WrapHandler)
}
