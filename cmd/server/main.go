package main

import (
	"context"
	"fmt"
	"log"
	"metrics/internal/server/config"
	"metrics/internal/server/logger"
	"metrics/internal/server/repository"
	"metrics/internal/server/router"

	"go.uber.org/zap"
)

func main() {
	loggerZap, err := logger.InitilazerLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	if err := run(loggerZap); err != nil {
		loggerZap.Fatal("Application failed", zap.Error(err))
	}
}

func run(loggerZap *zap.SugaredLogger) error {
	cfg, err := config.ParseFlags()
	if err != nil {
		loggerZap.Info(err.Error(), "failed to parse flags")
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	ctx := context.Background()
	rep, err := repository.NewMetricStorage(ctx, cfg, loggerZap)
	if err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	if err := router.Run(cfg, rep, loggerZap); err != nil {
		return fmt.Errorf("failed to run router: %w", err)
	}

	return nil
}
