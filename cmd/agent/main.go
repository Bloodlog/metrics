package main

import (
	"fmt"
	"log"
	"metrics/internal/config/agent"
	"metrics/internal/handlers"
	"metrics/internal/logger"
	repository2 "metrics/internal/repository"
	"metrics/internal/service"

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

	memoryRepository := repository2.NewMemoryRepository()
	systemRepository := repository2.NewSystemRepository()

	client := service.NewClient(configs.Address, configs.Key, configs.CryptoKey, loggerZap)

	applicationHandlers := handlers.NewAgentHandler(
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
