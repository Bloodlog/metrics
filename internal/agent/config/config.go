package config

import (
	"errors"
	"flag"
	"fmt"
	"metrics/internal/agent/dto"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

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

	flagKey        = "k"
	envKey         = "KEY"
	keyDescription = "Agent adds a HashSHA256 header with the computed hash"

	flagRateLimit        = "l"
	envRateLimit         = "RATE_LIMIT"
	rateLimitDescription = "Rate limit"

	flagCryptoKey        = "crypto-key"
	envCryptoKey         = "CRYPTO_KEY"
	cryptoKeyDescription = "Cryptographic encryption key"
)

func ParseFlags() (*dto.Config, error) {
	addressFlag := flag.String("a", defaultAddress, addressFlagDescription)
	reportIntervalFlag := flag.Int("r", defaultReportInterval, reportIntervalFlagDescription)
	pollIntervalFlag := flag.Int("p", defaultPollInterval, pollIntervalFlagDescription)
	keyFlag := flag.String(flagKey, "", keyDescription)
	rateLimitFlag := flag.Int(flagRateLimit, 1, rateLimitDescription)
	cryptoFlag := flag.String(flagCryptoKey, "", cryptoKeyDescription)
	configShort := flag.String("c", "", "Path to config file (short)")
	configLong := flag.String("config", "", "Path to config file (long)")
	flag.Parse()

	uknownArguments := flag.Args()
	if err := validateUnknownArgs(uknownArguments); err != nil {
		return nil, fmt.Errorf("read flags: %w", err)
	}

	return processFlags(
		*addressFlag,
		*reportIntervalFlag,
		*pollIntervalFlag,
		*keyFlag,
		*rateLimitFlag,
		*cryptoFlag,
		*configShort,
		*configLong,
	)
}

func processFlags(
	addressFlag string,
	reportIntervalFlag,
	pollIntervalFlag int,
	keyFlag string,
	rateLimitFlag int,
	cryptoKeyFlag string,
	configShort string,
	configLong string,
) (*dto.Config, error) {
	configPath := configLong
	if configPath == "" {
		configPath = configShort
	}
	if configPath != "" {
		if fromEnv, ok := os.LookupEnv("CONFIG"); ok {
			configPath = fromEnv
		}
	}
	var fileCfg dto.Config
	if configPath != "" {
		loaded, err := loadConfigFromFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("load config from file: %w", err)
		}
		fileCfg = *loaded
	}

	finalAddress, err := getStringValue(addressFlag, envAddress, fileCfg.Address)
	if err != nil {
		return nil, fmt.Errorf("read flag: %w", err)
	}

	host, port, err := parseAddress(finalAddress)
	if err != nil {
		return nil, fmt.Errorf("read flag address: %w", err)
	}
	address := "http://" + net.JoinHostPort(host, port)

	reportInterval, err := getIntValue(reportIntervalFlag, envReportInterval, fileCfg.ReportInterval)
	if err != nil {
		return nil, fmt.Errorf("read flag report interval: %w", err)
	}

	poolInterval, err := getIntValue(pollIntervalFlag, envPollInterval, fileCfg.PollInterval)
	if err != nil {
		return nil, fmt.Errorf("read flag pool interval: %w", err)
	}

	key, err := getStringValue(keyFlag, envKey, fileCfg.Key)
	if err != nil {
		key = ""
	}

	rateLimit, err := getIntValue(rateLimitFlag, envRateLimit, fileCfg.RateLimit)
	if err != nil {
		rateLimit = 1
	}

	cryptoKey, err := getStringValue(cryptoKeyFlag, envCryptoKey, fileCfg.CryptoKey)
	if err != nil {
		cryptoKey = ""
	}

	return &dto.Config{
		Address:        address,
		ReportInterval: reportInterval,
		PollInterval:   poolInterval,
		Batch:          false,
		Key:            key,
		RateLimit:      rateLimit,
		CryptoKey:      cryptoKey,
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

func getStringValue(flagValue string, envKey, fileVal string) (string, error) {
	if envValue, exists := os.LookupEnv(envKey); exists {
		return envValue, nil
	}

	if flagValue != "" {
		return flagValue, nil
	}

	if fileVal != "" {
		return fileVal, nil
	}

	return "", fmt.Errorf("missing required configuration: %s or flag value", envKey)
}

func getIntValue(flagValue int, envKey string, fileVal int) (int, error) {
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

	if fileVal != 0 {
		return fileVal, nil
	}

	return 0, fmt.Errorf("missing required configuration: %s or flag value", envKey)
}
