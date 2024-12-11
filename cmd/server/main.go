package main

import (
	"fmt"
	"log"
	"metrics/internal/server/config"
	"metrics/internal/server/repository"
	"metrics/internal/server/router"
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
	storage := repository.NewMemStorage()

	if err := router.Run(configs, storage); err != nil {
		return fmt.Errorf("failed to run router with provided configs and storage: %w", err)
	}

	return nil
}
