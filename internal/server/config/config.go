package config

import (
	"errors"
	"flag"
	"fmt"
	"metrics/internal/server/dto"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	flagHTTPAddress        = "a"
	envHTTPAddress         = "ADDRESS"
	defaultHTTPAddress     = "http://localhost:8080"
	descriptionHTTPAddress = "HTTP server address in the format host:port (default: localhost:8080)"
)

const (
	flagStoreInterval        = "i"
	envStoreInterval         = "STORE_INTERVAL"
	defaultStoreInterval     = 300
	descriptionStoreInterval = "interval fo store server"
)

const (
	flagRestore        = "r"
	envRestore         = "RESTORE"
	defaultRestore     = true
	descriptionRestore = "загружать ранее сохранённые значения из указанного файла при старте сервера"
)

const (
	flagFileStoragePath        = "f"
	envFileStoragePath         = "FILE_STORAGE_PATH"
	defaultFileStoragePath     = "metrics.json"
	descriptionFileStoragePath = "путь до файла, куда сохраняются текущие значения"
)

const (
	flagDatabaseDSN        = "d"
	envDatabaseDSN         = "DATABASE_DSN"
	defaultDatabaseDSN     = ""
	descriptionDatabaseDSN = "example postgres://username:password@localhost:5432/database_name"
)
const (
	flagKey        = "k"
	envKey         = "KEY"
	defaultKey     = ""
	descriptionKey = "Agent adds a HashSHA256 header with the computed hash"
)

const (
	flagCryptoKey        = "crypto-key"
	envCryptoKey         = "CRYPTO_KEY"
	cryptoKeyDescription = "Cryptographic encryption key"
)

func ParseFlags() (*dto.Config, error) {
	addressFlag := flag.String(flagHTTPAddress, defaultHTTPAddress, descriptionHTTPAddress)
	storeIntervalFlag := flag.Int(flagStoreInterval, defaultStoreInterval, descriptionStoreInterval)
	storagePathFlag := flag.String(flagFileStoragePath, defaultFileStoragePath, descriptionFileStoragePath)
	restoreFlag := flag.Bool(flagRestore, defaultRestore, descriptionRestore)
	addressDatabaseFlag := flag.String(flagDatabaseDSN, defaultDatabaseDSN, descriptionDatabaseDSN)
	keyFlag := flag.String(flagKey, defaultKey, descriptionKey)
	cryptoFlag := flag.String(flagCryptoKey, "", cryptoKeyDescription)
	enablePprof := flag.Bool("pprof", false, "enable pprof for debugging")
	configShort := flag.String("c", "", "Path to config file (short)")
	configLong := flag.String("config", "", "Path to config file (long)")
	flag.Parse()

	uknownArguments := flag.Args()
	if err := validateUnknownArgs(uknownArguments); err != nil {
		return nil, fmt.Errorf("read flag UnknownArgs: %w", err)
	}

	return processFlags(
		*addressFlag,
		*storeIntervalFlag,
		*storagePathFlag,
		*restoreFlag,
		*addressDatabaseFlag,
		*keyFlag,
		*cryptoFlag,
		*enablePprof,
		*configShort,
		*configLong,
	)
}

func processFlags(
	addressFlag string,
	storeIntervalFlag int,
	storagePathFlag string,
	restoreFlag bool,
	addressDatabaseFlag string,
	keyFlag string,
	cryptoKeyFlag string,
	enablePprof bool,
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

	finalAddress, err := getStringValue(addressFlag, envHTTPAddress, fileCfg.Address)
	if err != nil {
		finalAddress = ""
	}

	host, port, err := parseAddress(finalAddress)
	if err != nil {
		return nil, fmt.Errorf("read flag address: %w", err)
	}
	address := net.JoinHostPort(host, port)

	storeInterval, err := getIntValue(storeIntervalFlag, envStoreInterval, fileCfg.StoreInterval)
	if err != nil {
		return nil, fmt.Errorf("read flag report interval: %w", err)
	}

	storagePath, err := getStringValue(storagePathFlag, envFileStoragePath, fileCfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("read flag storage: %w", err)
	}

	restore := fileCfg.Restore
	restore, err = getBoolValue(restoreFlag, envRestore)
	if err != nil {
		return nil, fmt.Errorf("read flag restore: %w", err)
	}

	databaseDsn, err := getStringValue(addressDatabaseFlag, envDatabaseDSN, fileCfg.DatabaseDsn)
	if err != nil {
		databaseDsn = ""
	}

	key, err := getStringValue(keyFlag, envKey, fileCfg.Key)
	if err != nil {
		key = ""
	}

	cryptoKey, err := getStringValue(cryptoKeyFlag, envCryptoKey, fileCfg.CryptoKey)
	if err != nil {
		cryptoKey = ""
	}

	return &dto.Config{
		Address:         address,
		StoreInterval:   storeInterval,
		FileStoragePath: storagePath,
		DatabaseDsn:     databaseDsn,
		Restore:         restore,
		Key:             key,
		Debug:           enablePprof,
		CryptoKey:       cryptoKey,
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

func getStringValue(flagValue, envKey string, fileVal string) (string, error) {
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
