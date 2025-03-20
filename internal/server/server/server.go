package server

import (
	"fmt"
	"metrics/internal/server/config"
	"metrics/internal/server/repository"
	"metrics/internal/server/router"
	"net"
	"net/http"

	"go.uber.org/zap"
)

func ConfigureServerHandler(
	memStorage repository.MetricStorage,
	cfg *config.Config,
	logger *zap.SugaredLogger,
) error {
	handlerLogger := logger.With("r", "r")
	serverAddr := net.JoinHostPort(cfg.NetAddress.Host, cfg.NetAddress.Port)

	r := router.ConfigureServerHandler(memStorage, cfg, logger)
	handlerLogger.Infow(
		"Starting server",
		"addr", serverAddr,
	)
	err := http.ListenAndServe(serverAddr, r)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
