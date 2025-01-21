package repository

import (
	"context"
	"fmt"
	"metrics/internal/server/config"

	"go.uber.org/zap"
)

type FileRetryStorageWrapper struct {
	fileStorage *FileStorageWrapper
	cfg         *config.Config
	logger      *zap.SugaredLogger
}

func NewRetryFileStorage(
	ctx context.Context,
	cfg *config.Config,
	logger *zap.SugaredLogger,
) (*FileRetryStorageWrapper, error) {
	handlerLogger := logger.With("file-retry", "NewRetryFileStorageWrapper")

	var fileRepo *FileStorageWrapper
	err := retry(ctx, func() error {
		var err error
		fileRepo, err = NewFileStorageWrapper(ctx, cfg, logger)
		if err != nil {
			return fmt.Errorf("failed to create FileStorageWrapper: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create FileStorageWrapper: %w", err)
	}

	fileRetryrepo := &FileRetryStorageWrapper{
		fileStorage: fileRepo,
		cfg:         cfg,
		logger:      handlerLogger,
	}

	handlerLogger.Infof("Using file storage with retry: %s", cfg.FileStoragePath)
	return fileRetryrepo, nil
}

func (fr *FileRetryStorageWrapper) SetGauge(ctx context.Context, name string, value float64) error {
	return retry(ctx, func() error {
		err := fr.fileStorage.SetGauge(ctx, name, value)
		if err != nil {
			return &RetriableError{Err: err}
		}
		return nil
	})
}

func (fr *FileRetryStorageWrapper) GetGauge(ctx context.Context, name string) (float64, error) {
	var result float64
	err := retry(ctx, func() error {
		var err error
		result, err = fr.fileStorage.GetGauge(ctx, name)
		if err != nil {
			return &RetriableError{Err: err}
		}
		return nil
	})
	return result, err
}

func (fr *FileRetryStorageWrapper) SetCounter(ctx context.Context, name string, value uint64) error {
	return retry(ctx, func() error {
		err := fr.fileStorage.SetCounter(ctx, name, value)
		if err != nil {
			return &RetriableError{Err: err}
		}
		return nil
	})
}

func (fr *FileRetryStorageWrapper) GetCounter(ctx context.Context, name string) (uint64, error) {
	var result uint64
	err := retry(ctx, func() error {
		var err error
		result, err = fr.fileStorage.GetCounter(ctx, name)
		if err != nil {
			return &RetriableError{Err: err}
		}
		return nil
	})
	return result, err
}

func (fr *FileRetryStorageWrapper) Gauges(ctx context.Context) (map[string]float64, error) {
	var result map[string]float64
	err := retry(ctx, func() error {
		var err error
		result, err = fr.fileStorage.Gauges(ctx)
		if err != nil {
			return &RetriableError{Err: err}
		}
		return nil
	})
	return result, err
}

func (fr *FileRetryStorageWrapper) Counters(ctx context.Context) (map[string]uint64, error) {
	var result map[string]uint64
	err := retry(ctx, func() error {
		var err error
		result, err = fr.fileStorage.Counters(ctx)
		if err != nil {
			return &RetriableError{Err: err}
		}
		return nil
	})
	return result, err
}

func (fr *FileRetryStorageWrapper) UpdateCounterAndGauges(
	ctx context.Context,
	name string,
	value uint64,
	gauges map[string]float64,
) error {
	return retry(ctx, func() error {
		err := fr.fileStorage.UpdateCounterAndGauges(ctx, name, value, gauges)
		if err != nil {
			return &RetriableError{Err: err}
		}
		return nil
	})
}