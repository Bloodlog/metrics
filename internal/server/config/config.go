package config

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type Config struct {
	NetAddress NetAddress
	Debug      bool
}

type NetAddress struct {
	Host string
	Port int
}

func ParseFlags(flagAddress string, unknownArgs []string, envAddress string) (*Config, error) {
	if err := validateUnknownArgs(unknownArgs); err != nil {
		return nil, err
	}

	finalAddress, err := getFinalAddress(flagAddress, envAddress)
	if err != nil {
		return nil, err
	}

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

func getFinalAddress(flagValue string, envVar string) (string, error) {
	if envVar != "" {
		return envVar, nil
	}

	if flagValue != "" {
		return flagValue, nil
	}

	return "", errors.New("no address provided via flag or environment variable")
}

func parseAddress(address string) (string, int, error) {
	parsedURL, err := url.Parse(address)
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse address: %w", err)
	}

	host := parsedURL.Hostname()
	portStr := parsedURL.Port()

	if host == "" || portStr == "" {
		return "", 0, fmt.Errorf("invalid address format: %s (expected host:port)", address)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port value: %w", err)
	}

	return host, port, nil
}
