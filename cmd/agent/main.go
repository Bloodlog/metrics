package main

import (
	"fmt"
	"log"
	"metrics/internal/agent/config"
	"metrics/internal/agent/handlers"
	"metrics/internal/agent/repository"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	configs, err := config.ParseFlags()
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	storage := repository.NewRepository()

	if err := handlers.Handle(configs, storage); err != nil {
		return fmt.Errorf("failed to handle configs and storage: %w", err)
	}

	return nil
}
