package repository

import (
	"context"
	"fmt"
	"metrics/internal/server/dto"

	"go.uber.org/zap"
)

type RetryDBRepository struct {
	storage *DBRepository
	cfg     *dto.Config
}

func NewRetryBRepository(
	ctx context.Context,
	cfg *dto.Config,
	logger *zap.SugaredLogger,
) (*RetryDBRepository, error) {
	handlerLogger := logger.With("retry", "NewRetryBRepository")
	handlerLogger.Info("Attempting to connect to the database...")

	var db *DBRepository
	err := retry(ctx, func() error {
		var err error
		db, err = NewDBRepository(ctx, cfg, logger)
		if err != nil {
			return fmt.Errorf("failed to initialize DBRepository: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize DBRepository after retries: %w", err)
	}

	dbRetryRepository := &RetryDBRepository{
		storage: db,
		cfg:     cfg,
	}

	handlerLogger.Info("RetryDBRepository initialized successfully.")
	return dbRetryRepository, nil
}

func (r *RetryDBRepository) SetGauge(ctx context.Context, name string, value float64) (float64, error) {
	var result float64
	err := retry(ctx, func() error {
		var err error
		result, err = r.storage.SetGauge(ctx, name, value)
		return err
	})
	return result, err
}

func (r *RetryDBRepository) GetGauge(ctx context.Context, name string) (float64, error) {
	var result float64
	err := retry(ctx, func() error {
		var err error
		result, err = r.storage.GetGauge(ctx, name)
		return err
	})
	return result, err
}

func (r *RetryDBRepository) SetCounter(ctx context.Context, name string, value uint64) (uint64, error) {
	var result uint64
	err := retry(ctx, func() error {
		var err error
		result, err = r.storage.SetCounter(ctx, name, value)
		return err
	})
	return result, err
}

func (r *RetryDBRepository) GetCounter(ctx context.Context, name string) (uint64, error) {
	var result uint64
	err := retry(ctx, func() error {
		var err error
		result, err = r.storage.GetCounter(ctx, name)
		return err
	})
	return result, err
}

func (r *RetryDBRepository) Gauges(ctx context.Context) (map[string]float64, error) {
	var result map[string]float64
	err := retry(ctx, func() error {
		var err error
		result, err = r.storage.Gauges(ctx)
		return err
	})
	return result, err
}

func (r *RetryDBRepository) Counters(ctx context.Context) (map[string]uint64, error) {
	var result map[string]uint64
	err := retry(ctx, func() error {
		var err error
		result, err = r.storage.Counters(ctx)
		return err
	})
	return result, err
}

func (r *RetryDBRepository) UpdateCounterAndGauges(
	ctx context.Context,
	counters map[string]uint64,
	gauges map[string]float64,
) error {
	return retry(ctx, func() error {
		return r.storage.UpdateCounterAndGauges(ctx, counters, gauges)
	})
}
