package repository

import (
	"context"
	"fmt"
	"sync"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]uint64
	mu       *sync.RWMutex
}

func NewMemStorage() (MetricStorage, error) {
	memStorage := &MemStorage{
		mu:       &sync.RWMutex{},
		gauges:   make(map[string]float64),
		counters: make(map[string]uint64),
	}

	return memStorage, nil
}

func (ms *MemStorage) SetGauge(ctx context.Context, name string, value float64) (float64, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.gauges[name] = value

	return ms.gauges[name], nil
}

func (ms *MemStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	value, exists := ms.gauges[name]
	if !exists {
		return 0, fmt.Errorf("gauge metric '%s' not found", name)
	}
	return value, nil
}

func (ms *MemStorage) SetCounter(ctx context.Context, name string, value uint64) (uint64, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.counters[name] += value

	return ms.counters[name], nil
}

func (ms *MemStorage) GetCounter(ctx context.Context, name string) (uint64, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	value, exists := ms.counters[name]
	if !exists {
		return 0, fmt.Errorf("counter metric '%s' not found", name)
	}
	return value, nil
}

func (ms *MemStorage) Gauges(ctx context.Context) (map[string]float64, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	result := make(map[string]float64, len(ms.gauges))
	for k, v := range ms.gauges {
		result[k] = v
	}
	return result, nil
}

func (ms *MemStorage) Counters(ctx context.Context) (map[string]uint64, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	result := make(map[string]uint64, len(ms.counters))
	for k, v := range ms.counters {
		result[k] = v
	}
	return result, nil
}

func (ms *MemStorage) UpdateCounterAndGauges(
	ctx context.Context,
	counters map[string]uint64,
	gauges map[string]float64,
) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	for counterName, counterValue := range counters {
		_, err := ms.SetCounter(ctx, counterName, counterValue)
		if err != nil {
			return fmt.Errorf("error saving counter: %w", err)
		}
	}

	for gaugeName, gaugeValue := range gauges {
		_, err := ms.SetGauge(ctx, gaugeName, gaugeValue)
		if err != nil {
			return fmt.Errorf("error saving metrics: %w", err)
		}
	}

	return nil
}
