package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	NetAddress NetAddress
	Debug      bool
}

type NetAddress struct {
	Host string
	Port int
}

func ParseFlags() (*Config, error) {
	const DefaultAddress = "localhost:8080"
	address := flag.String("a", DefaultAddress, "HTTP server address in the format host:port (default: localhost:8080)")

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

	return &Config{
		NetAddress: NetAddress{Host: host, Port: port},
		Debug:      false,
	}, nil
}
