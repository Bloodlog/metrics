package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
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

func ParseFlags(flagAddress string, flagReportInterval int, flagPollInterval int, unknownArgs []string) (*Config, error) {
	if err := validateUnknownArgs(unknownArgs); err != nil {
		return nil, err
	}

	const EnvAddress = "ADDRESS"
	const EnvReportInterval = "REPORT_INTERVAL"
	const EnvPoolInterval = "POLL_INTERVAL"
	envAddress := os.Getenv(EnvAddress)
	envReportInterval := os.Getenv(EnvReportInterval)
	envPollInterval := os.Getenv(EnvPoolInterval)

	finalAddress, err := getFinalAddress(flagAddress, envAddress)
	if err != nil {
		return nil, err
	}

	host, port, err := parseAddress(finalAddress)
	if err != nil {
		return nil, err
	}

	reportInterval, err := getInterval(flagReportInterval, envReportInterval)
	if err != nil {
		return nil, fmt.Errorf("failed to parse REPORT_INTERVAL: %w", err)
	}

	pollInterval, err := getInterval(flagPollInterval, envPollInterval)
	if err != nil {
		return nil, fmt.Errorf("failed to parse POLL_INTERVAL: %w", err)
	}

	return &Config{
		NetAddress:     NetAddress{Host: host, Port: port},
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
		Debug:          false,
	}, nil
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

func getInterval(flagValue int, envVar string) (int, error) {
	if envVar != "" {
		interval, err := strconv.Atoi(envVar)
		if err != nil {
			return 0, err
		}

		return interval, nil
	}

	return flagValue, nil
}

func validateUnknownArgs(unknownArgs []string) error {
	if len(unknownArgs) > 0 {
		return fmt.Errorf("error: unknown flags or arguments detected: %v", unknownArgs)
	}
	return nil
}
