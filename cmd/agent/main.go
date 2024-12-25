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
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	logger, err := getLogger()
	if err != nil {
		return fmt.Errorf("logger fail: %w", err)
	}
	logger.Info("Logger initialized successfully")

	configs, err := config.ParseFlags()
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	logger.Info("Flags and env parsed")

	storage := repository.NewRepository()
	logger.Info("Repository initialized successfully")

	if err := handlers.Handle(configs, storage, logger); err != nil {
		return fmt.Errorf("failed to handle configs and storage: %w", err)
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
