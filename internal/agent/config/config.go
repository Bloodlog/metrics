package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

type Config struct {
	NetAddress     NetAddress
	ReportInterval int
	PollInterval   int
	Debug          bool
}

func ParseFlags() (*Config, error) {
	address := flag.String("a", "localhost:8080", "HTTP server address in the format host:port (default: localhost:8080)")
	reportIntervalArg := flag.Int("r", 10, "Overrides the metric reporting frequency to the server (default: 10 seconds)")
	pollIntervalArg := flag.Int("p", 2, "Overrides the metric polling frequency from the runtime package (default: 2 seconds)")

	flag.Parse()

	if len(flag.Args()) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Error: unknown flags detected")
		os.Exit(1)
	}

	parts := strings.Split(*address, ":")
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		parts = strings.Split(envRunAddr, ":")
	}
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid address format: %s (expected host:port)", *address)
	}

	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to convert port to number: %w", err)
	}

	reportInterval := *reportIntervalArg
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		reportInterval, _ = strconv.Atoi(envReportInterval)
	}

	pollInterval := *pollIntervalArg
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		pollInterval, _ = strconv.Atoi(envPollInterval)
	}

	return &Config{
		NetAddress:     NetAddress{Host: host, Port: port},
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
		Debug:          false,
	}, nil
}
