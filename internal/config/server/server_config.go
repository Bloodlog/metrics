package server

import (
	"flag"
	"fmt"
	"metrics/internal/config"
	"net"
	"os"
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

func ParseFlags() (*config.ServerConfig, error) {
	addressFlag := flag.String(flagHTTPAddress, defaultHTTPAddress, descriptionHTTPAddress)
	storeIntervalFlag := flag.Int(flagStoreInterval, defaultStoreInterval, descriptionStoreInterval)
	storagePathFlag := flag.String(flagFileStoragePath, defaultFileStoragePath, descriptionFileStoragePath)
	restoreFlag := flag.Bool(flagRestore, defaultRestore, descriptionRestore)
	addressDatabaseFlag := flag.String(flagDatabaseDSN, defaultDatabaseDSN, descriptionDatabaseDSN)
	keyFlag := flag.String(flagKey, defaultKey, descriptionKey)
	cryptoFlag := flag.String(flagCryptoKey, "", cryptoKeyDescription)
	enablePprof := flag.Bool("pprof", false, "enable pprof for debugging")
	trustedSubnet := flag.String("t", "", "CIDR")
	configShort := flag.String("c", "", "Path to config file (short)")
	configLong := flag.String("config", "", "Path to config file (long)")
	flag.Parse()

	uknownArguments := flag.Args()
	if err := config.ValidateUnknownArgs(uknownArguments); err != nil {
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
		*trustedSubnet,
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
	trustedSubnetFlag string,
	configShort string,
	configLong string,
) (*config.ServerConfig, error) {
	configPath := configLong
	if configPath == "" {
		configPath = configShort
	}
	if configPath != "" {
		if fromEnv, ok := os.LookupEnv("CONFIG"); ok {
			configPath = fromEnv
		}
	}
	var fileCfg config.ServerConfig
	if configPath != "" {
		err := config.LoadConfigFromFile(configPath, &fileCfg)
		if err != nil {
			return nil, fmt.Errorf("load config from file: %w", err)
		}
	}

	finalAddress, err := config.GetStringValue(addressFlag, envHTTPAddress, fileCfg.Address)
	if err != nil {
		finalAddress = ""
	}

	host, port, err := config.ParseAddress(finalAddress)
	if err != nil {
		return nil, fmt.Errorf("read flag address: %w", err)
	}
	address := net.JoinHostPort(host, port)

	storeInterval, err := config.GetIntValue(storeIntervalFlag, envStoreInterval, fileCfg.StoreInterval)
	if err != nil {
		return nil, fmt.Errorf("read flag report interval: %w", err)
	}

	storagePath, err := config.GetStringValue(storagePathFlag, envFileStoragePath, fileCfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("read flag storage: %w", err)
	}

	restore, err := config.GetBoolValue(restoreFlag, envRestore)
	if err != nil {
		return nil, fmt.Errorf("read flag restore: %w", err)
	}

	databaseDsn, err := config.GetStringValue(addressDatabaseFlag, envDatabaseDSN, fileCfg.DatabaseDsn)
	if err != nil {
		databaseDsn = ""
	}

	key, err := config.GetStringValue(keyFlag, envKey, fileCfg.Key)
	if err != nil {
		key = ""
	}

	cryptoKey, err := config.GetStringValue(cryptoKeyFlag, envCryptoKey, fileCfg.CryptoKey)
	if err != nil {
		cryptoKey = ""
	}

	trustedSubnet, err := config.GetStringValue(trustedSubnetFlag, "TRUSTED_SUBNET", fileCfg.TrustedSubnet)
	if err != nil {
		trustedSubnet = ""
	}

	var trustedNet *net.IPNet
	if trustedSubnet != "" {
		ip, cidr, err := net.ParseCIDR(trustedSubnet)
		if err != nil {
			return nil, fmt.Errorf("invalid trusted_subnet: %w", err)
		}
		if ip != nil && cidr != nil {
			trustedNet = cidr
		}
	}

	return &config.ServerConfig{
		Address:         address,
		StoreInterval:   storeInterval,
		FileStoragePath: storagePath,
		DatabaseDsn:     databaseDsn,
		Restore:         restore,
		Key:             key,
		Debug:           enablePprof,
		CryptoKey:       cryptoKey,
		TrustedSubnet:   trustedSubnet,
		TrustedNet:      trustedNet,
	}, nil
}
