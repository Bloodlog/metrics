package main

import (
	"flag"
	"fmt"
	"log"
	"metrics/internal/server/config"
	"metrics/internal/server/repository"
	"metrics/internal/server/router"
	"os"
)

const DefaultAddress = "http://localhost:8080"
const EnvAddress = "ADDRESS"
const DefaultAddressDescription = "HTTP server address in the format host:port (default: localhost:8080)"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	addressFlag := flag.String("a", DefaultAddress, DefaultAddressDescription)
	flag.Parse()

	envAddress := os.Getenv(EnvAddress)

	configs, err := config.ParseFlags(*addressFlag, flag.Args(), envAddress)
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	storage := repository.NewMemStorage()

	if err := router.Run(configs, storage); err != nil {
		return fmt.Errorf("failed to run router with provided configs and storage: %w", err)
	}

	return nil
}
