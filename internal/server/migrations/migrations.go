package migrations

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(ctx context.Context, conn *pgxpool.Pool, logger *zap.SugaredLogger) error {
	handlerLogger := logger.With("package", "migrations")
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = 'metrics'
		)
	`
	var exists bool
	if err := conn.QueryRow(ctx, query).Scan(&exists); err != nil {
		return fmt.Errorf("error check exist table: %w", err)
	}

	if exists {
		handlerLogger.Info("Table metrics already exist")
		return nil
	}

	createTableQuery := `
		CREATE TABLE metrics (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			value DOUBLE PRECISION NOT NULL
		)
	`
	if _, err := conn.Exec(ctx, createTableQuery); err != nil {
		return fmt.Errorf("error create table: %w", err)
	}

	handlerLogger.Info("Migration create table metrics success")
	return nil
}
