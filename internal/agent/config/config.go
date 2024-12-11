package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port string
}

type Config struct {
	NetAddress     NetAddress
	ReportInterval int
	PollInterval   int
	Debug          bool
}

const (
	DefaultAddress        = "http://localhost:8080"
	DefaultReportInterval = 10
	DefaultPollInterval   = 2

	EnvAddress        = "ADDRESS"
	EnvReportInterval = "REPORT_INTERVAL"
	EnvPollInterval   = "POLL_INTERVAL"

	AddressFlagDescription        = "HTTP server address in the format host:port (default: localhost:8080)"
	ReportIntervalFlagDescription = "Overrides the metric reporting frequency to the server (default: 10 seconds)"
	PollIntervalFlagDescription   = "Overrides the metric polling frequency (default: 2 seconds)"
)

func ParseFlags() (*Config, error) {
	addressFlag := flag.String("a", DefaultAddress, AddressFlagDescription)
	reportIntervalFlag := flag.Int("r", DefaultReportInterval, ReportIntervalFlagDescription)
	pollIntervalFlag := flag.Int("p", DefaultPollInterval, PollIntervalFlagDescription)
	flag.Parse()

	uknownArguments := flag.Args()
	if err := validateUnknownArgs(uknownArguments); err != nil {
		log.Printf("error: unknown flags or arguments detected: %v", uknownArguments)
		return nil, err
	}

	finalAddress, err := getStringValue(*addressFlag, EnvAddress)
	if err != nil {
		log.Printf("error: invalid address: %v", err)
		return nil, err
	}

	host, port, err := parseAddress(finalAddress)
	if err != nil {
		log.Printf("error: invalid address: %v", err)
		return nil, err
	}

	reportInterval, err := getIntValue(*reportIntervalFlag, EnvReportInterval)
	if err != nil {
		log.Printf("Warning: invalid value for %s", EnvReportInterval)
		return nil, err
	}

	poolInterval, err := getIntValue(*pollIntervalFlag, EnvPollInterval)
	if err != nil {
		log.Printf("Warning: invalid value for %s", EnvPollInterval)
		return nil, err
	}

	return &Config{
		NetAddress:     NetAddress{Host: host, Port: port},
		ReportInterval: reportInterval,
		PollInterval:   poolInterval,
		Debug:          false,
	}, nil
}

func parseAddress(address string) (string, string, error) {
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = "http://" + address
	}
	parsedURL, err := url.Parse(address)
	if err != nil {
		return "", "", errors.New("failed to parse address (expected host:port)")
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()

	if host == "" || port == "" {
		return "", "", errors.New("failed to parse address (expected host:port)")
	}

	return host, port, nil
}

func validateUnknownArgs(unknownArgs []string) error {
	if len(unknownArgs) > 0 {
		return errors.New("unknown flags or arguments detected")
	}
	return nil
}

func getStringValue(flagValue, envKey string) (string, error) {
	if envValue, exists := os.LookupEnv(envKey); exists {
		return envValue, nil
	}

	if flagValue != "" {
		return flagValue, nil
	}

	return "", fmt.Errorf("missing required configuration: %s or flag value", envKey)
}

func getIntValue(flagValue int, envKey string) (int, error) {
	if envValue, exists := os.LookupEnv(envKey); exists {
		parsedValue, err := strconv.Atoi(envValue)
		if err != nil {
			return 0, fmt.Errorf("invalid value for environment variable %s: %s", envKey, envValue)
		}
		return parsedValue, nil
	}

	if flagValue != 0 {
		return flagValue, nil
	}

	return 0, fmt.Errorf("missing required configuration: %s or flag value", envKey)
}
