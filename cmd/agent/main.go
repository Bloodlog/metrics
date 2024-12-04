package main

import (
	"flag"
	"fmt"
	"log"
	"metrics/internal/agent/config"
	"metrics/internal/agent/handlers"
	"metrics/internal/agent/repository"
	"os"
)

const DefaultAddress = "http://localhost:8080"
const DefaultReportInterval = 10
const DefaultPoolInterval = 2
const EnvAddress = "ADDRESS"
const EnvReportInterval = "REPORT_INTERVAL"
const EnvPoolInterval = "POLL_INTERVAL"
const AddressFlagDescription = "HTTP server address in the format host:port (default: localhost:8080)"
const DefaultReportIntervalDescription = "Overrides the metric reporting frequency to the server (default: 10 seconds)"
const DefaultPoolIntervalDescription = "Overrides the metric polling frequency (default: 2 seconds)"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	addressFlag := flag.String("a", DefaultAddress, AddressFlagDescription)
	reportIntervalFlag := flag.Int("r", DefaultReportInterval, DefaultReportIntervalDescription)
	pollIntervalFlag := flag.Int("p", DefaultPoolInterval, DefaultPoolIntervalDescription)
	flag.Parse()

	envAddress := os.Getenv(EnvAddress)
	envReportInterval := os.Getenv(EnvReportInterval)
	envPollInterval := os.Getenv(EnvPoolInterval)

	configs, err := config.ParseFlags(
		*addressFlag,
		*reportIntervalFlag,
		*pollIntervalFlag,
		flag.Args(),
		envAddress,
		envReportInterval,
		envPollInterval,
	)
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	storage := repository.NewRepository()

	if err := handlers.Handle(configs, storage); err != nil {
		return fmt.Errorf("failed to handle configs and storage: %w", err)
	}

	return nil
}
