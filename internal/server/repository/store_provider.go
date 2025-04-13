package repository

import (
	"context"
	"errors"
	"fmt"
	"metrics/internal/server/dto"
	"time"

	"go.uber.org/zap"
)

type MetricType string

type RetriableError struct {
	Err error
}

func (e *RetriableError) Error() string {
	return e.Err.Error()
}

func (e *RetriableError) Unwrap() error {
	return e.Err
}

func NewMetricStorage(ctx context.Context, cfg *dto.Config, logger *zap.SugaredLogger) (MetricStorage, error) {
	storageType := resolve(cfg)

	switch storageType {
	case "memory":
		storage, _ := NewMemStorage()

		return storage, nil
	case "file":
		storage, _ := NewFileStorageWrapper(ctx, cfg, logger)

		return storage, nil
	case "file-retry":
		storage, err := NewRetryFileStorage(ctx, cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create NewRetryFileStorage: %w", err)
		}

		return storage, nil
	case "database":
		storage, err := NewDBRepository(ctx, cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create NewDBRepository: %w", err)
		}

		return storage, nil
	case "database-retry":
		storage, err := NewRetryBRepository(ctx, cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create RetryDBRepository: %w", err)
		}

		return storage, nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
}

func resolve(cfg *dto.Config) string {
	if cfg.DatabaseDsn != "" {
		return "database-retry"
	}
	if cfg.FileStoragePath != "" {
		return "file-retry"
	}

	return "memory"
}

func retry(ctx context.Context, operation func() error) error {
	retryIntervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	var lastError error
	lastError = operation()
	if lastError == nil {
		return nil
	}

	var retriableErr *RetriableError
	if !errors.As(lastError, &retriableErr) {
		return fmt.Errorf("operation failed: %w", lastError)
	}

	for i, interval := range retryIntervals {
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation canceled: %w", ctx.Err())
		default:
			if i > 0 {
				time.Sleep(interval)
			}

			lastError = operation()
			if lastError == nil {
				return nil
			}

			if !errors.As(lastError, &retriableErr) {
				return fmt.Errorf("operation failed: %w", lastError)
			}
		}
	}

	return fmt.Errorf("operation failed after retries: %w", lastError)
}
