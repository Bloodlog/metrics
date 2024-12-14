package main

import (
	"fmt"
	"log"
	"metrics/internal/server/config"
	"metrics/internal/server/repository"
	"metrics/internal/server/router"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("logger failed: %w", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("failed to sync logger: %v\n", err)
		}
	}()
	sugar = *logger.Sugar()

	configs, err := config.ParseFlags(sugar)
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	storage := repository.NewMemStorage()

	if err := router.Run(configs, storage, sugar); err != nil {
		return fmt.Errorf("failed to run router with provided configs and storage: %w", err)
	}

	return nil
}
