package main

import (
	"fmt"
	"log"
	"metrics/internal/agent/config"
	"metrics/internal/agent/handlers"
	"metrics/internal/agent/repository"

	"go.uber.org/zap"
)

func main() {
	logger, err := getLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	if err := run(logger); err != nil {
		logger.Fatal("Application failed", zap.Error(err))
	}
}

func run(logger *zap.SugaredLogger) error {
	configs, err := config.ParseFlags()
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	storage := repository.NewRepository()

	applicationHandlers := handlers.NewHandlers(configs, storage, logger)
	if err := applicationHandlers.Handle(); err != nil {
		return fmt.Errorf("application failed: %w", err)
	}

	return nil
}

func getLogger() (*zap.SugaredLogger, error) {
	configLogger := zap.NewDevelopmentConfig()
	configLogger.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	logger, err := configLogger.Build()
	if err != nil {
		return nil, fmt.Errorf("logger initialization failed: %w", err)
	}

	return logger.Sugar(), err
}
