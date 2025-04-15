package router

import (
	"crypto/rsa"
	"log"
	"metrics/internal/config"
	"metrics/internal/handlers/api"
	"metrics/internal/handlers/web"
	middleware2 "metrics/internal/middleware"
	"metrics/internal/repository"
	"metrics/internal/security"
	"metrics/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "metrics/swagger"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func ConfigureServerHandler(
	memStorage repository.MetricStorage,
	cfg *config.ServerConfig,
	logger *zap.SugaredLogger,
) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware2.LoggingMiddleware(logger),
		middleware2.DecompressionMiddleware(logger),
		middleware2.CheckHashMiddleware(logger, cfg.Key),
		middleware2.ResponseHashMiddleware(cfg.Key),
		middleware2.ResponseCompressionMiddleware(logger),
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
	cfg *config.ServerConfig,
	memStorage repository.MetricStorage,
	logger *zap.SugaredLogger,
) {
	metricService := service.NewMetricService(memStorage, logger)
	apiHandler := api.NewHandler(metricService, logger)
	webHandler := web.NewHandler(metricService, logger)

	var privateKey *rsa.PrivateKey
	if cfg.CryptoKey != "" {
		var err error
		privateKey, err = security.LoadRSAPrivateKeyFromFile(cfg.CryptoKey)
		if err != nil {
			log.Fatalf("failed to load private key: %v", err)
		}
	}

	r.Route("/updates", func(r chi.Router) {
		r.Use(middleware2.DecryptMiddleware(privateKey))
		r.Post("/", apiHandler.UpdatesHandler())
	})
	r.Route("/update", func(r chi.Router) {
		r.Use(middleware2.DecryptMiddleware(privateKey))
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
