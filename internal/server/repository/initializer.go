package repository

import (
	"context"
	"fmt"
	"metrics/internal/server/migrations"

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
	if databaseDsn != "" {
		pool, err := pgxpool.New(ctx, databaseDsn)
		if err != nil {
			return nil, fmt.Errorf("failed to create pool: %w", err)
		}

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
			if err := fileRepo.LoadFromFile(ctx); err != nil {
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
