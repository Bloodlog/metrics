package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

const (
	DefaultAddress         = "http://localhost:8080"
	EnvAddress             = "ADDRESS"
	AddressFlagDescription = "HTTP server address in the format host:port (default: localhost:8080)"
)

type Config struct {
	NetAddress NetAddress
	Debug      bool
}

type NetAddress struct {
	Host string
	Port string
}

func ParseFlags() (*Config, error) {
	addressFlag := flag.String("a", DefaultAddress, AddressFlagDescription)
	flag.Parse()

	if err := validateUnknownArgs(flag.Args()); err != nil {
		return nil, err
	}

	finalAddress := getStringValue(*addressFlag, EnvAddress, DefaultAddress)

	host, port, err := parseAddress(finalAddress)
	if err != nil {
		return nil, err
	}

	return &Config{
		NetAddress: NetAddress{Host: host, Port: port},
		Debug:      false,
	}, nil
}

func validateUnknownArgs(unknownArgs []string) error {
	if len(unknownArgs) > 0 {
		return fmt.Errorf("error: unknown flags or arguments detected: %v", unknownArgs)
	}
	return nil
}

func parseAddress(address string) (string, string, error) {
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = "http://" + address
	}
	parsedURL, err := url.Parse(address)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse address: %w", err)
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()

	if host == "" || port == "" {
		return "", "", fmt.Errorf("invalid address format: %s (expected host:port)", address)
	}

	return host, port, nil
}

func getStringValue(flagValue, envKey, defaultValue string) string {
	if flagValue != defaultValue {
		return flagValue
	}
	if envValue, exists := os.LookupEnv(envKey); exists {
		return envValue
	}
	return defaultValue
}
