package main

import (
	"context"
	"fmt"
	"log"
	server2 "metrics/internal/config/server"
	"metrics/internal/logger"
	"metrics/internal/repository"
	"metrics/internal/server"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var (
	version     = "N/A"
	buildTime   = "N/A"
	buildCommit = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", buildTime)
	fmt.Printf("Build commit: %s\n", buildCommit)
	loggerZap, err := logger.InitilazerLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	if err = run(loggerZap); err != nil {
		loggerZap.Fatal("Application failed", zap.Error(err))
	}
}

// @title Metrics API
// @version 1.0
// @description API для управления метриками
// @host 127.0.0.1:8080
// @BasePath /.
func run(loggerZap *zap.SugaredLogger) error {
	cfg, err := server2.ParseFlags()
	if err != nil {
		loggerZap.Info(err.Error(), "failed to parse flags")
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	memStorage, err := repository.NewMetricStorage(ctx, cfg, loggerZap)
	if err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	pprofServer := server.InitPprof(cfg, loggerZap)

	httpServer, err := server.ConfigureServerHandler(memStorage, cfg, loggerZap)
	if err != nil {
		return fmt.Errorf("failed to run router: %w", err)
	}

	<-ctx.Done()
	loggerZap.Info("Shutdown signal received")
	if err = httpServer.Shutdown(context.Background()); err != nil {
		loggerZap.Info("HTTP server Shutdown: %v", err)
	}

	if pprofServer != nil {
		if err = pprofServer.Shutdown(context.Background()); err != nil {
			loggerZap.Errorw("PProf server Shutdown failed", "error", err)
		} else {
			loggerZap.Info("PProf server gracefully stopped")
		}
	}

	return nil
}
