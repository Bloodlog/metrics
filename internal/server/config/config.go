package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

const (
	AddressFlag            = "a"
	EnvAddress             = "ADDRESS"
	DefaultAddress         = "http://localhost:8080"
	AddressFlagDescription = "HTTP server address in the format host:port (default: localhost:8080)"
)

const (
	StoreIntervalFlg         = "i"
	EnvStoreInterval         = "STORE_INTERVAL"
	DefaultStoreInterval     = 300
	StoreIntervalDescription = "Interval fo store server"
)

const (
	RestoreFlag        = "r"
	EnvRestore         = "RESTORE"
	DefaultRestore     = true
	RestoreDescription = "загружать ранее сохранённые значения из указанного файла при старте сервера"
)

const (
	FileStoragePathFlag        = "f"
	EnvFileStoragePath         = "FILE_STORAGE_PATH"
	DefaultFileStoragePath     = "metrics.json"
	FileStoragePathDescription = "путь до файла, куда сохраняются текущие значения"
)

const (
	nameError = "config"
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

func ParseFlags(logger zap.SugaredLogger) (*Config, error) {
	addressFlag := flag.String(AddressFlag, DefaultAddress, AddressFlagDescription)
	storeIntervalFlag := flag.Int(StoreIntervalFlg, DefaultStoreInterval, StoreIntervalDescription)
	storagePathFlag := flag.String(FileStoragePathFlag, DefaultFileStoragePath, FileStoragePathDescription)
	restoreFlag := flag.Bool(RestoreFlag, DefaultRestore, RestoreDescription)

	flag.Parse()

	uknownArguments := flag.Args()
	if err := validateUnknownArgs(uknownArguments); err != nil {
		logger.Infoln(err.Error(), nameError, "read flag UnknownArgs")
		return nil, err
	}

	finalAddress, err := getStringValue(*addressFlag, EnvAddress)
	if err != nil {
		logger.Infoln(err.Error(), nameError, "read flag address")
		return nil, err
	}

	host, port, err := parseAddress(finalAddress)
	if err != nil {
		logger.Infoln(err.Error(), nameError, "read flag address")
		return nil, err
	}

	storeInterval, err := getIntValue(*storeIntervalFlag, EnvStoreInterval)
	if err != nil {
		logger.Infoln(err.Error(), nameError, "read flag report interval")
		return nil, err
	}

	storagePath, err := getStringValue(*storagePathFlag, EnvFileStoragePath)
	if err != nil {
		logger.Infoln(err.Error(), nameError, "read flag storage")
		return nil, err
	}

	restore, err := getBoolValue(*restoreFlag, EnvRestore)
	if err != nil {
		logger.Infoln(err.Error(), nameError, "read flag restore")
		return nil, err
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
