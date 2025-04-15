package main

import (
	"fmt"
	"log"
	"metrics/internal/agent/clients"
	"metrics/internal/agent/handlers"
	"metrics/internal/agent/repository"
	"metrics/internal/config/agent"
	"metrics/internal/logger"

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

	if err := run(loggerZap); err != nil {
		loggerZap.Fatal("Application failed", zap.Error(err))
	}
}

func run(loggerZap *zap.SugaredLogger) error {
	configs, err := agent.ParseAgentFlags()
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	memoryRepository := repository.NewMemoryRepository()
	systemRepository := repository.NewSystemRepository()

	client := clients.NewClient(configs.Address, configs.Key, configs.CryptoKey, loggerZap)

	applicationHandlers := handlers.NewHandlers(
		client,
		configs,
		memoryRepository,
		systemRepository,
		loggerZap,
	)
	if err = applicationHandlers.Handle(); err != nil {
		return fmt.Errorf("application failed: %w", err)
	}

	return nil
}
