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

const (
	addressFlag            = "a"
	envAddress             = "ADDRESS"
	defaultAddress         = "http://localhost:8080"
	addressFlagDescription = "HTTP server address in the format host:port (default: localhost:8080)"
)

const (
	storeIntervalFlg         = "i"
	envStoreInterval         = "STORE_INTERVAL"
	defaultStoreInterval     = 300
	storeIntervalDescription = "interval fo store server"
)

const (
	restoreFlag        = "r"
	envRestore         = "RESTORE"
	defaultRestore     = true
	restoreDescription = "загружать ранее сохранённые значения из указанного файла при старте сервера"
)

const (
	fileStoragePathFlag        = "f"
	envFileStoragePath         = "FILE_STORAGE_PATH"
	defaultFileStoragePath     = "metrics.json"
	fileStoragePathDescription = "путь до файла, куда сохраняются текущие значения"
)

type Config struct {
	NetAddress      NetAddress
	FileStoragePath string
	StoreInterval   int
	Restore         bool
}

type NetAddress struct {
	Host string
	Port string
}

func ParseFlags() (*Config, error) {
	addressFlag := flag.String(addressFlag, defaultAddress, addressFlagDescription)
	storeIntervalFlag := flag.Int(storeIntervalFlg, defaultStoreInterval, storeIntervalDescription)
	storagePathFlag := flag.String(fileStoragePathFlag, defaultFileStoragePath, fileStoragePathDescription)
	restoreFlag := flag.Bool(restoreFlag, defaultRestore, restoreDescription)

	flag.Parse()

	uknownArguments := flag.Args()
	if err := validateUnknownArgs(uknownArguments); err != nil {
		return nil, fmt.Errorf("read flag UnknownArgs: %w", err)
	}

	finalAddress, err := getStringValue(*addressFlag, envAddress)
	if err != nil {
		return nil, fmt.Errorf("read flag address: %w", err)
	}

	host, port, err := parseAddress(finalAddress)
	if err != nil {
		return nil, fmt.Errorf("read flag address: %w", err)
	}

	storeInterval, err := getIntValue(*storeIntervalFlag, envStoreInterval)
	if err != nil {
		return nil, fmt.Errorf("read flag report interval: %w", err)
	}

	storagePath, err := getStringValue(*storagePathFlag, envFileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("read flag storage: %w", err)
	}

	restore, err := getBoolValue(*restoreFlag, envRestore)
	if err != nil {
		return nil, fmt.Errorf("read flag restore: %w", err)
	}

	return &Config{
		NetAddress:      NetAddress{Host: host, Port: port},
		StoreInterval:   storeInterval,
		FileStoragePath: storagePath,
		Restore:         restore,
	}, nil
}

func validateUnknownArgs(unknownArgs []string) error {
	if len(unknownArgs) > 0 {
		return fmt.Errorf("unknown flags or arguments detected: %v", unknownArgs)
	}
	return nil
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

func getBoolValue(flagValue bool, envKey string) (bool, error) {
	if envValue, exists := os.LookupEnv(envKey); exists {
		parsedValue, err := strconv.ParseBool(envValue)
		if err != nil {
			return false, fmt.Errorf("invalid boolean value for environment variable %s: %s", envKey, envValue)
		}
		return parsedValue, nil
	}

	return flagValue, nil
}
