package repository

import (
	"context"
	"fmt"
	"metrics/internal/server/repository/migrations"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"metrics/internal/server/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBRepository struct {
	pool   *pgxpool.Pool
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func NewDBRepository(
	ctx context.Context,
	cfg *config.Config,
	logger *zap.SugaredLogger,
) (*DBRepository, error) {
	handlerLogger := logger.With("database", "NewDBRepository")
	pool, err := initPool(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}

	repository := &DBRepository{pool: pool, cfg: cfg, logger: handlerLogger}

	err = migrations.Migrate(ctx, pool, logger)
	if err != nil {
		return nil, fmt.Errorf("error migrate: %w", err)
	}
	handlerLogger.Infoln("Sucess migrate")

	return repository, nil
}

func initPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseDsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the DSN: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, &RetriableError{Err: err}
	}
	return pool, nil
}

func (r *DBRepository) SetGauge(ctx context.Context, name string, value float64) (float64, error) {
	query := `
		INSERT INTO metrics (name, value, mtype)
		VALUES ($1, $2, 'gauge')
		ON CONFLICT (name) DO UPDATE SET value = $2, mtype = 'gauge'
		RETURNING value
	`
	var newValue float64
	err := r.pool.QueryRow(ctx, query, name, value).Scan(&newValue)
	if err != nil {
		return 0, fmt.Errorf("error setting gauge '%s': %w", name, err)
	}

	return newValue, nil
}

func (r *DBRepository) GetGauge(ctx context.Context, name string) (float64, error) {
	query := `SELECT value FROM metrics WHERE name = $1 AND mtype = 'gauge'`
	var value float64
	err := r.pool.QueryRow(ctx, query, name).Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("error getting gauge '%s': %w", name, err)
	}
	return value, nil
}

func (r *DBRepository) SetCounter(ctx context.Context, name string, value uint64) (uint64, error) {
	query := `
		INSERT INTO metrics (name, delta, mtype)
		VALUES ($1, $2, 'counter')
		ON CONFLICT (name) DO UPDATE SET delta = metrics.delta + $2, mtype = 'counter'
		RETURNING delta
	`
	var newValue uint64
	err := r.pool.QueryRow(ctx, query, name, value).Scan(&newValue)
	if err != nil {
		return 0, fmt.Errorf("error setting gauge '%s': %w", name, err)
	}

	return newValue, nil
}

func (r *DBRepository) GetCounter(ctx context.Context, name string) (uint64, error) {
	query := `SELECT delta FROM metrics WHERE name = $1 AND mtype = 'counter'`
	var value uint64
	err := r.pool.QueryRow(ctx, query, name).Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("error getting counter '%s': %w", name, err)
	}
	return value, nil
}

func (r *DBRepository) Gauges(ctx context.Context) (map[string]float64, error) {
	query := `SELECT name, value FROM metrics WHERE mtype = 'gauge'`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error get Gauges: %w", err)
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
	return gauges, nil
}

func (r *DBRepository) Counters(ctx context.Context) (map[string]uint64, error) {
	query := `SELECT name, delta FROM metrics WHERE mtype = 'counter'`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error get counters: %w", err)
	}
	defer rows.Close()

	counters := make(map[string]uint64)
	for rows.Next() {
		var name string
		var delta float64
		if err := rows.Scan(&name, &delta); err != nil {
			fmt.Printf("error scanning row: %v\n", err)
			continue
		}
		counters[name] = uint64(delta)
	}
	return counters, nil
}

func (r *DBRepository) UpdateCounterAndGauges(
	ctx context.Context,
	counters map[string]uint64,
	gauges map[string]float64,
) error {
	batch := new(pgx.Batch)

	upsertCounterQuery := `
		INSERT INTO metrics (name, delta, mtype)
		VALUES ($1, $2, 'counter')
		ON CONFLICT (name) DO UPDATE SET delta = metrics.delta + $2, mtype = 'counter'`

	for counterName, counterValue := range counters {
		batch.Queue(upsertCounterQuery, counterName, counterValue)
	}

	upsertGaugesQuery := `
		INSERT INTO metrics (name, value, mtype)
		VALUES ($1, $2, 'gauge')
		ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value, mtype = 'gauge'`
	for gaugeName, gaugeValue := range gauges {
		batch.Queue(upsertGaugesQuery, gaugeName, gaugeValue)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer func(results pgx.BatchResults) {
		err := results.Close()
		if err != nil {
			r.logger.Infoln("Error send batch", err)
		}
	}(results)

	return nil
}
