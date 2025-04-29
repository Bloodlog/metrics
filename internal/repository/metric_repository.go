package repository

import "context"

type MetricStorage interface {
	SetGauge(ctx context.Context, name string, value float64) (float64, error)
	GetGauge(ctx context.Context, name string) (float64, error)
	SetCounter(ctx context.Context, name string, value uint64) (uint64, error)
	GetCounter(ctx context.Context, name string) (uint64, error)
	Gauges(ctx context.Context) (map[string]float64, error)
	Counters(ctx context.Context) (map[string]uint64, error)
	UpdateCounterAndGauges(ctx context.Context, counters map[string]uint64, gauges map[string]float64) error
	Shutdown(ctx context.Context)
}
