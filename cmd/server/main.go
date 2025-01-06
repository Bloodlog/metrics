package main

import (
	"fmt"
	"log"
	"metrics/internal/server/config"
	"metrics/internal/server/repository"
	"metrics/internal/server/router"

	"go.uber.org/zap"
)

func main() {
	logger, err := getLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	if err := run(logger); err != nil {
		logger.Fatal("Application failed", zap.Error(err))
	}
}

func run(logger *zap.SugaredLogger) error {
	configs, err := config.ParseFlags()
	if err != nil {
		logger.Info(err.Error(), "failed to parse flags")
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	storage, err := getRepository(configs.FileStoragePath, configs.StoreInterval, configs.Restore, logger)
	if err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	if configs.StoreInterval > 0 && configs.FileStoragePath != "" {
		go func() {
			if err := storage.AutoSave(); err != nil {
				logger.Info("Error during AutoSave", "error", err)
			}
		}()
		logger.Info("AutoSave initialized")
	}

	if err := router.Run(configs, storage, logger); err != nil {
		return fmt.Errorf("failed to run router: %w", err)
	}

	return nil
}

func getRepository(
	fileStoragePath string,
	storeInterval int,
	restore bool,
	logger *zap.SugaredLogger,
) (repository.MetricStorage, error) {
	logger.Info("Initializing repository...")

	storage := repository.NewMemStorage()
	logger.Info("Memory storage initialized")

	if fileStoragePath != "" {
		logger.Infof("Using file storage: %s", fileStoragePath)
		fileStorage := repository.NewFileStorageWrapper(storage, fileStoragePath, storeInterval)

		if restore {
			logger.Info("Restore is enabled, loading from file...")
			if err := fileStorage.LoadFromFile(); err != nil {
				return nil, fmt.Errorf("error loading metrics: %w", err)
			}
		}

		logger.Info("Using in-memory repository")
		return fileStorage, nil
	}

	logger.Info("Repository use memory")
	return storage, nil
}

func getLogger() (*zap.SugaredLogger, error) {
	configLogger := zap.NewDevelopmentConfig()
	configLogger.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	logger, err := configLogger.Build()
	if err != nil {
		return nil, fmt.Errorf("logger initialization failed: %w", err)
	}

	return logger.Sugar(), err
}
