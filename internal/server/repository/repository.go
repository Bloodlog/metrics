package repository

import (
	"context"
	"errors"
	"fmt"
	"metrics/internal/server/config"
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

func IsRetriableError(err error) bool {
	var retriableErr *RetriableError
	return errors.As(err, &retriableErr)
}

type MetricStorage interface {
	SetGauge(ctx context.Context, name string, value float64) (float64, error)
	GetGauge(ctx context.Context, name string) (float64, error)
	SetCounter(ctx context.Context, name string, value uint64) (uint64, error)
	GetCounter(ctx context.Context, name string) (uint64, error)
	Gauges(ctx context.Context) (map[string]float64, error)
	Counters(ctx context.Context) (map[string]uint64, error)
	UpdateCounterAndGauges(ctx context.Context, name string, value uint64, gauges map[string]float64) error
}

func NewMetricStorage(ctx context.Context, cfg *config.Config, logger *zap.SugaredLogger) (MetricStorage, error) {
	storageType := resolve(cfg)

	switch storageType {
	case "memory":
		storage, _ := NewMemStorage(ctx)

		return storage, nil
	case "file":
		storage, _ := NewFileStorageWrapper(ctx, cfg, logger)

		return storage, nil
	case "file-retry":
		storage, _ := NewRetryFileStorage(ctx, cfg, logger)

		return storage, nil
	case "database":
		storage, _ := NewDBRepository(ctx, cfg, logger)

		return storage, nil
	case "database-retry":
		storage, _ := NewRetryBRepository(ctx, cfg, logger)

		return storage, nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
}

func resolve(cfg *config.Config) string {
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
