package server

import (
	"metrics/internal/config"
	"metrics/internal/repository"
	"metrics/internal/router"
	"net/http"

	"go.uber.org/zap"
)

func ConfigureServerHandler(
	memStorage repository.MetricStorage,
	cfg *config.ServerConfig,
	logger *zap.SugaredLogger,
) error {
	handlerLogger := logger.With("r", "r")

	r := router.ConfigureServerHandler(memStorage, cfg, logger)
	handlerLogger.Infow(
		"Starting server",
		"addr", cfg.Address,
	)
	go func() {
		err := http.ListenAndServe(cfg.Address, r)
		if err != nil {
			handlerLogger.Info("Failed to start server", "err", err)
		}
	}()
	return nil
}
