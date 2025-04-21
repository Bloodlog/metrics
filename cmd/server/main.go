package main

import (
	"context"
	"fmt"
	"log"
	"metrics/internal/config"
	server2 "metrics/internal/config/server"
	"metrics/internal/logger"
	"metrics/internal/repository"
	"metrics/internal/server"
	"net/http"
	_ "net/http/pprof"
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

	initPprof(cfg, loggerZap)
	if err = server.ConfigureServerHandler(memStorage, cfg, loggerZap); err != nil {
		return fmt.Errorf("failed to run router: %w", err)
	}

	<-ctx.Done()
	loggerZap.Info("Shutdown signal received")

	return nil
}

func initPprof(cfg *config.ServerConfig, zapLog *zap.SugaredLogger) {
	if cfg.Debug {
		go func() {
			err := http.ListenAndServe("0.0.0.0"+":6060", nil)
			if err != nil {
				zapLog.Info(err.Error(), "failed start profiler")
			}
		}()
	}
}
