package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBRepository struct {
	pool *pgxpool.Pool
}

func NewDBRepository(pool *pgxpool.Pool) *DBRepository {
	return &DBRepository{pool: pool}
}

func (r *DBRepository) SetGauge(ctx context.Context, name string, value float64) error {
	var exists bool
	checkQuery := `SELECT EXISTS (SELECT 1 FROM metrics WHERE name = $1)`
	err := r.pool.QueryRow(ctx, checkQuery, name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking existence of gauge '%s': %w", name, err)
	}

	if exists {
		updateQuery := `UPDATE metrics SET value = $2 WHERE name = $1`
		_, err = r.pool.Exec(ctx, updateQuery, name, value)
		if err != nil {
			return fmt.Errorf("error updating gauge '%s': %w", name, err)
		}
	} else {
		insertQuery := `INSERT INTO metrics (name, value) VALUES ($1, $2)`
		_, err = r.pool.Exec(ctx, insertQuery, name, value)
		if err != nil {
			return fmt.Errorf("error inserting gauge '%s': %w", name, err)
		}
	}

	return nil
}

func (r *DBRepository) GetGauge(ctx context.Context, name string) (float64, error) {
	query := `SELECT value FROM metrics WHERE name = $1`
	var value float64
	err := r.pool.QueryRow(ctx, query, name).Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("error getting gauge '%s': %w", name, err)
	}
	return value, nil
}

func (r *DBRepository) SetCounter(ctx context.Context, name string, value uint64) error {
	var exists bool
	checkQuery := `SELECT EXISTS (SELECT 1 FROM metrics WHERE name = $1)`
	err := r.pool.QueryRow(ctx, checkQuery, name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking existence of counter '%s': %w", name, err)
	}

	if exists {
		updateQuery := `UPDATE metrics SET value = value + $2 WHERE name = $1`
		_, err = r.pool.Exec(ctx, updateQuery, name, float64(value))
		if err != nil {
			return fmt.Errorf("error updating counter '%s': %w", name, err)
		}
	} else {
		insertQuery := `INSERT INTO metrics (name, value) VALUES ($1, $2)`
		_, err = r.pool.Exec(ctx, insertQuery, name, float64(value))
		if err != nil {
			return fmt.Errorf("error inserting counter '%s': %w", name, err)
		}
	}

	return nil
}

func (r *DBRepository) GetCounter(ctx context.Context, name string) (uint64, error) {
	query := `SELECT value FROM metrics WHERE name = $1`
	var value float64
	err := r.pool.QueryRow(ctx, query, name).Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("error getting counter '%s': %w", name, err)
	}
	return uint64(value), nil
}

func (r *DBRepository) Gauges(ctx context.Context) map[string]float64 {
	query := `SELECT name, value FROM metrics`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	gauges := make(map[string]float64)
	for rows.Next() {
		var name string
		var value float64
		if err := rows.Scan(&name, &value); err != nil {
			continue
		}
		gauges[name] = value
	}
	return gauges
}

func (r *DBRepository) Counters(ctx context.Context) map[string]uint64 {
	query := `SELECT name, value FROM metrics`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	counters := make(map[string]uint64)
	for rows.Next() {
		var name string
		var value float64
		if err := rows.Scan(&name, &value); err != nil {
			fmt.Printf("error scanning row: %v\n", err)
			continue
		}
		counters[name] = uint64(value)
	}
	return counters
}

func (r *DBRepository) AutoSave(ctx context.Context) error {
	return nil
}

func (r *DBRepository) LoadFromFile(ctx context.Context) error {
	return nil
}

func (r *DBRepository) SaveToFile(ctx context.Context) error {
	return nil
}

func (r *DBRepository) WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			err = fmt.Errorf("panic recovered during transaction: %v", p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	err = fn(tx)
	if err != nil {
		return fmt.Errorf("transaction function execution failed: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
