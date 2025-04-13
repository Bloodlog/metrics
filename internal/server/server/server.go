package server

import (
	"fmt"
	"metrics/internal/server/dto"
	"metrics/internal/server/repository"
	"metrics/internal/server/router"
	"net/http"

	"go.uber.org/zap"
)

func ConfigureServerHandler(
	memStorage repository.MetricStorage,
	cfg *dto.Config,
	logger *zap.SugaredLogger,
) error {
	handlerLogger := logger.With("r", "r")

	r := router.ConfigureServerHandler(memStorage, cfg, logger)
	handlerLogger.Infow(
		"Starting server",
		"addr", cfg.Address,
	)
	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
