package config

import (
	"errors"
	"flag"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	ErrParseAddress   = errors.New("failed to parse address (expected host:port)")
	ErrArgumentsCount = errors.New("unknown flags or arguments detected")
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

	if err := validateUnknownArgs(flag.Args()); err != nil {
		return nil, err
	}

	finalAddress := getStringValue(*addressFlag, EnvAddress, DefaultAddress)

	host, port, err := parseAddress(finalAddress)
	if err != nil {
		return nil, err
	}

	return &Config{
		NetAddress:     NetAddress{Host: host, Port: port},
		ReportInterval: getIntValue(*reportIntervalFlag, EnvReportInterval, DefaultReportInterval),
		PollInterval:   getIntValue(*pollIntervalFlag, EnvPollInterval, DefaultPollInterval),
		Debug:          false,
	}, nil
}

func parseAddress(address string) (string, string, error) {
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = "http://" + address
	}
	parsedURL, err := url.Parse(address)
	if err != nil {
		return "", "", ErrParseAddress
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()

	if host == "" || port == "" {
		return "", "", ErrParseAddress
	}

	return host, port, nil
}

func validateUnknownArgs(unknownArgs []string) error {
	if len(unknownArgs) > 0 {
		log.Printf("error: unknown flags or arguments detected: %v", unknownArgs)
		return ErrArgumentsCount
	}
	return nil
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

func getIntValue(flagValue int, envKey string, defaultValue int) int {
	if flagValue != defaultValue {
		return flagValue
	}
	if envValue, exists := os.LookupEnv(envKey); exists {
		if parsedValue, err := strconv.Atoi(envValue); err == nil {
			return parsedValue
		}
		log.Printf("Warning: invalid value for %s, using default: %d", envKey, defaultValue)
	}
	return defaultValue
}
