package agent

import (
	"flag"
	"fmt"
	"metrics/internal/config"
	"net"
	"os"
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

	flagAgentKey   = "k"
	envAgentKey    = "KEY"
	keyDescription = "Agent adds a HashSHA256 header with the computed hash"

	flagRateLimit        = "l"
	envRateLimit         = "RATE_LIMIT"
	rateLimitDescription = "Rate limit"

	flagAgentCryptoKey        = "crypto-key"
	envAgentCryptoKey         = "CRYPTO_KEY"
	cryptoAgentKeyDescription = "Cryptographic encryption key"
)

func ParseAgentFlags() (*config.AgentConfig, error) {
	addressFlag := flag.String("a", defaultAddress, addressFlagDescription)
	reportIntervalFlag := flag.Int("r", defaultReportInterval, reportIntervalFlagDescription)
	pollIntervalFlag := flag.Int("p", defaultPollInterval, pollIntervalFlagDescription)
	keyFlag := flag.String(flagAgentKey, "", keyDescription)
	rateLimitFlag := flag.Int(flagRateLimit, 1, rateLimitDescription)
	cryptoFlag := flag.String(flagAgentCryptoKey, "", cryptoAgentKeyDescription)
	configShort := flag.String("c", "", "Path to config file (short)")
	configLong := flag.String("config", "", "Path to config file (long)")
	flag.Parse()

	uknownArguments := flag.Args()
	if err := config.ValidateUnknownArgs(uknownArguments); err != nil {
		return nil, fmt.Errorf("read flags: %w", err)
	}

	return processAgentFlags(
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

func processAgentFlags(
	addressFlag string,
	reportIntervalFlag,
	pollIntervalFlag int,
	keyFlag string,
	rateLimitFlag int,
	cryptoKeyFlag string,
	configShort string,
	configLong string,
) (*config.AgentConfig, error) {
	configPath := configLong
	if configPath == "" {
		configPath = configShort
	}
	if configPath != "" {
		if fromEnv, ok := os.LookupEnv("CONFIG"); ok {
			configPath = fromEnv
		}
	}
	var fileCfg config.AgentConfig
	if configPath != "" {
		err := config.LoadConfigFromFile(configPath, &fileCfg)
		if err != nil {
			return nil, fmt.Errorf("load config from file: %w", err)
		}
	}

	finalAddress, err := config.GetStringValue(addressFlag, envAddress, fileCfg.Address)
	if err != nil {
		return nil, fmt.Errorf("read flag: %w", err)
	}

	host, port, err := config.ParseAddress(finalAddress)
	if err != nil {
		return nil, fmt.Errorf("read flag address: %w", err)
	}
	address := "http://" + net.JoinHostPort(host, port)

	reportInterval, err := config.GetIntValue(reportIntervalFlag, envReportInterval, fileCfg.ReportInterval)
	if err != nil {
		return nil, fmt.Errorf("read flag report interval: %w", err)
	}

	poolInterval, err := config.GetIntValue(pollIntervalFlag, envPollInterval, fileCfg.PollInterval)
	if err != nil {
		return nil, fmt.Errorf("read flag pool interval: %w", err)
	}

	key, err := config.GetStringValue(keyFlag, envAgentKey, fileCfg.Key)
	if err != nil {
		key = ""
	}

	rateLimit, err := config.GetIntValue(rateLimitFlag, envRateLimit, fileCfg.RateLimit)
	if err != nil {
		rateLimit = 1
	}

	cryptoKey, err := config.GetStringValue(cryptoKeyFlag, envAgentCryptoKey, fileCfg.CryptoKey)
	if err != nil {
		cryptoKey = ""
	}

	return &config.AgentConfig{
		Address:        address,
		ReportInterval: reportInterval,
		PollInterval:   poolInterval,
		Batch:          true,
		Key:            key,
		RateLimit:      rateLimit,
		CryptoKey:      cryptoKey,
		Grpc:           true,
	}, nil
}
