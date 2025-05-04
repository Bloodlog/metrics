package server

import (
	"fmt"
	"metrics/internal/config"
	"metrics/internal/repository"
	"metrics/internal/router"
	"net/http"

	_ "net/http/pprof"

	"go.uber.org/zap"
)

func ConfigureServerHandler(
	memStorage repository.MetricStorage,
	cfg *config.ServerConfig,
	logger *zap.SugaredLogger,
) (*http.Server, error) {
	handlerLogger := logger.With("r", "r")

	r := router.ConfigureServerHandler(memStorage, cfg, logger)
	handlerLogger.Infow(
		"Starting server",
		"addr", cfg.Address,
	)
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil {
		return srv, fmt.Errorf("listen and server has failed: %w", err)
	}

	return srv, nil
}

func InitPprof() (*http.Server, error) {
	pprofServer := &http.Server{
		Addr: "0.0.0.0:6060",
	}
	if err := pprofServer.ListenAndServe(); err != nil {
		return nil, fmt.Errorf("listen and server has failed: %w", err)
	}

	return pprofServer, nil
}
