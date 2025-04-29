package server

import (
	"errors"
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
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			handlerLogger.Errorw("Failed to start server", "err", err)
		}
	}()

	return srv, nil
}

func InitPprof(cfg *config.ServerConfig, zapLog *zap.SugaredLogger) *http.Server {
	if cfg.Debug {
		return nil
	}
	pprofServer := &http.Server{
		Addr: "0.0.0.0:6060",
	}
	go func() {
		if err := pprofServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zapLog.Warnw("failed to start profiler", "err", err)
		}
	}()

	return pprofServer
}
