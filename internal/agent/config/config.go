package config

import (
	"errors"
	"flag"
	"fmt"
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
	Batch          bool
	ClientSetting  ClientSetting
}

type ClientSetting struct {
	MaxNumberAttempts  int
	RetryWaitSecond    int
	RetryMaxWaitSecond int
}

const (
	defaultAddress        = "http://localhost:8080"
	defaultReportInterval = 10
	defaultPollInterval   = 2

	envAddress        = "ADDRESS"
	envReportInterval = "REPORT_INTERVAL"
	envPollInterval   = "POLL_INTERVAL"

	addressFlagDescription        = "HTTP server address in the format host:port (default: localhost:8080)"
	reportIntervalFlagDescription = "Overrides the metric reporting frequency to the server (default: 10 seconds)"
	pollIntervalFlagDescription   = "Overrides the metric polling frequency (default: 2 seconds)"
)

const (
	maxNumberAttempts  = 3
	retryWaitSecond    = 2
	retryMaxWaitSecond = 5
)

func ParseFlags() (*Config, error) {
	addressFlag := flag.String("a", defaultAddress, addressFlagDescription)
	reportIntervalFlag := flag.Int("r", defaultReportInterval, reportIntervalFlagDescription)
	pollIntervalFlag := flag.Int("p", defaultPollInterval, pollIntervalFlagDescription)
	flag.Parse()

	uknownArguments := flag.Args()
	if err := validateUnknownArgs(uknownArguments); err != nil {
		return nil, fmt.Errorf("read flags: %w", err)
	}

	finalAddress, err := getStringValue(*addressFlag, envAddress)
	if err != nil {
		return nil, fmt.Errorf("read flag: %w", err)
	}

	host, port, err := parseAddress(finalAddress)
	if err != nil {
		return nil, fmt.Errorf("read flag address: %w", err)
	}

	reportInterval, err := getIntValue(*reportIntervalFlag, envReportInterval)
	if err != nil {
		return nil, fmt.Errorf("read flag report interval: %w", err)
	}

	poolInterval, err := getIntValue(*pollIntervalFlag, envPollInterval)
	if err != nil {
		return nil, fmt.Errorf("read flag pool interval: %w", err)
	}

	return &Config{
		NetAddress:     NetAddress{Host: host, Port: port},
		ReportInterval: reportInterval,
		PollInterval:   poolInterval,
		Batch:          false,
		ClientSetting: ClientSetting{
			MaxNumberAttempts:  maxNumberAttempts,
			RetryWaitSecond:    retryWaitSecond,
			RetryMaxWaitSecond: retryMaxWaitSecond,
		},
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
		return fmt.Errorf("unknown flags or arguments detected: %v", unknownArgs)
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
