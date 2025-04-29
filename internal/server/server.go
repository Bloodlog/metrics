package server

import (
	"errors"
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
		if errors.Is(err, http.ErrServerClosed) {
			return nil, nil
		}
		return nil, fmt.Errorf("listen and server has failed: %w", err)
	}

	return srv, nil
}

func InitPprof(cfg *config.ServerConfig, zapLog *zap.SugaredLogger) (*http.Server, error) {
	if cfg.Debug {
		return nil, nil
	}
	pprofServer := &http.Server{
		Addr: "0.0.0.0:6060",
	}
	if err := pprofServer.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return nil, nil
		}

		return nil, fmt.Errorf("listen and server has failed: %w", err)
	}

	return pprofServer, nil
}
