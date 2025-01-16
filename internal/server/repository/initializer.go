package repository

import (
	"context"
	"fmt"
	"metrics/internal/server/migrations"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func InitializeRepository(
	ctx context.Context,
	logger *zap.SugaredLogger,
	databaseDsn,
	filePath string,
	storeInterval int,
	restore bool) (MetricStorage, error) {
	handlerLogger := logger.With("package", "repository")
	const maxAttempts = 3
	if databaseDsn != "" {
		handlerLogger.Info("Attempting to connect to the database...")
		pool, err := connectDBWithRetry(ctx, databaseDsn, maxAttempts, handlerLogger)
		if err != nil {
			return nil, fmt.Errorf("failed to create pool: %w", err)
		}

		handlerLogger.Info("Running migrations...")
		err = migrations.Migrate(ctx, pool, logger)
		if err != nil {
			return nil, fmt.Errorf("error migrate: %w", err)
		}

		handlerLogger.Info("Using DB storage")
		return NewDBRepository(pool), nil
	}

	if filePath != "" {
		logger.Infof("Using file storage: %s", filePath)
		memRepo := NewMemStorage()
		fileRepo := NewFileStorageWrapper(memRepo, filePath, storeInterval)

		if restore {
			logger.Info("Restore is enabled, loading from file...")
			if err := loadFromFileWithRetry(ctx, fileRepo, maxAttempts, handlerLogger); err != nil {
				return nil, fmt.Errorf("error loading metrics: %w", err)
			}
		}

		if storeInterval > 0 {
			go func() {
				if err := fileRepo.AutoSave(ctx); err != nil {
					handlerLogger.Info("Error during AutoSave", "error", err)
				}
			}()
		}

		handlerLogger.Info("Using in-memory repository")
		return fileRepo, nil
	}

	logger.Info("Repository use memory")
	return NewMemStorage(), nil
}

func connectDBWithRetry(
	ctx context.Context,
	databaseDsn string,
	maxRetries int,
	logger *zap.SugaredLogger,
) (*pgxpool.Pool, error) {
	delays := []time.Duration{1 * time.Second, 3 * time.Second}
	var err error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		logger.Infof("Attempt %d to connect to the database...", attempt+1)
		pool, err := pgxpool.New(ctx, databaseDsn)
		if err == nil {
			if connErr := testDBConnection(ctx, pool); connErr == nil {
				logger.Info("Successfully connected to the database.")
				return pool, nil
			} else {
				logger.Warnw("Database connection test failed", "error", connErr)
				pool.Close()
			}
		}

		if attempt < maxRetries {
			time.Sleep(delays[attempt])
		}
	}

	return nil, fmt.Errorf("operation failed after %d attempts: %w", maxRetries+1, err)
}

func testDBConnection(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	conn.Release()
	return nil
}

func loadFromFileWithRetry(
	ctx context.Context,
	fileRepo *FileStorageWrapper,
	maxRetries int,
	logger *zap.SugaredLogger,
) error {
	delays := []time.Duration{1 * time.Second, 3 * time.Second}
	var err error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		logger.Infof("Attempt %d to load metrics from file...", attempt+1)
		err = fileRepo.LoadFromFile(ctx)
		if err == nil {
			logger.Info("Successfully loaded metrics from file.")
			return nil
		}

		logger.Warnw("Failed to load metrics from file", "attempt", attempt+1, "error", err)
		if attempt < maxRetries {
			time.Sleep(delays[attempt])
		}
	}

	return fmt.Errorf("operation failed after %d attempts: %w", maxRetries+1, err)
}
