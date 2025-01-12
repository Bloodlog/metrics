package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type MetricType string

type MetricStorage interface {
	SetGauge(ctx context.Context, name string, value float64) error
	GetGauge(ctx context.Context, name string) (float64, error)
	SetCounter(ctx context.Context, name string, value uint64) error
	GetCounter(ctx context.Context, name string) (uint64, error)
	Gauges(ctx context.Context) map[string]float64
	Counters(ctx context.Context) map[string]uint64
	AutoSave(ctx context.Context) error
	LoadFromFile(ctx context.Context) error
	SaveToFile(ctx context.Context) error
	WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error
}
