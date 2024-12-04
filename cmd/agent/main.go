package main

import (
	"flag"
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
	const DefaultAddress = "http://localhost:8080"
	const DefaultReportInterval = 10
	const DefaultPoolInterval = 2

	addressFlag := flag.String("a", DefaultAddress, "HTTP server address in the format host:port (default: localhost:8080)")
	reportIntervalFlag := flag.Int("r", DefaultReportInterval, "Overrides the metric reporting frequency to the server (default: 10 seconds)")
	pollIntervalFlag := flag.Int("p", DefaultPoolInterval, "Overrides the metric polling frequency from the runtime package (default: 2 seconds)")
	flag.Parse()

	configs, err := config.ParseFlags(*addressFlag, *reportIntervalFlag, *pollIntervalFlag, flag.Args())
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	storage := repository.NewRepository()

	return handlers.Handle(configs, storage)
}
