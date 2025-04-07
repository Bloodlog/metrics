package main

import (
	"context"
	"fmt"
	"log"
	"metrics/internal/server/config"
	"metrics/internal/server/logger"
	"metrics/internal/server/repository"
	"metrics/internal/server/server"

	"net/http"
	_ "net/http/pprof"

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
	cfg, err := config.ParseFlags()
	if err != nil {
		loggerZap.Info(err.Error(), "failed to parse flags")
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	ctx := context.Background()
	memStorage, err := repository.NewMetricStorage(ctx, cfg, loggerZap)
	if err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	initPprof(cfg, loggerZap)
	if err = server.ConfigureServerHandler(memStorage, cfg, loggerZap); err != nil {
		return fmt.Errorf("failed to run router: %w", err)
	}

	return nil
}

func initPprof(cfg *config.Config, zapLog *zap.SugaredLogger) {
	if cfg.Debug {
		go func() {
			err := http.ListenAndServe(cfg.NetAddress.Host+":6060", nil)
			if err != nil {
				zapLog.Info(err.Error(), "failed start profiler")
			}
		}()
	}
}
