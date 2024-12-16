package main

import (
	"fmt"
	"log"
	"metrics/internal/agent/config"
	"metrics/internal/agent/handlers"
	"metrics/internal/agent/repository"

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
	storage := repository.NewRepository()

	if err := handlers.Handle(configs, storage, sugar); err != nil {
		return fmt.Errorf("failed to handle configs and storage: %w", err)
	}

	return nil
}
