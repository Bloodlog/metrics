package migrations

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(ctx context.Context, conn *pgxpool.Pool, logger *zap.SugaredLogger) error {
	handlerLogger := logger.With("package", "migrations")

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS metrics (
			id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			name VARCHAR(200) NOT NULL UNIQUE,
			value DOUBLE PRECISION NULL,
			delta INTEGER NULL,
			mtype VARCHAR(200) NOT NULL
		)
	`

	if _, err := conn.Exec(ctx, createTableQuery); err != nil {
		return fmt.Errorf("error create table: %w", err)
	}

	handlerLogger.Info("Migration create table metrics success")
	return nil
}
