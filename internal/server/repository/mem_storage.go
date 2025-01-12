package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]uint64
	mu       *sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		mu:       &sync.RWMutex{},
		gauges:   make(map[string]float64),
		counters: make(map[string]uint64),
	}
}

func (ms *MemStorage) SetGauge(ctx context.Context, name string, value float64) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.gauges[name] = value

	return nil
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

func (ms *MemStorage) SetCounter(ctx context.Context, name string, value uint64) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.counters[name] += value

	return nil
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

func (ms *MemStorage) Gauges(ctx context.Context) map[string]float64 {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	result := make(map[string]float64, len(ms.gauges))
	for k, v := range ms.gauges {
		result[k] = v
	}
	return result
}

func (ms *MemStorage) Counters(ctx context.Context) map[string]uint64 {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	result := make(map[string]uint64, len(ms.counters))
	for k, v := range ms.counters {
		result[k] = v
	}
	return result
}

func (ms *MemStorage) AutoSave(ctx context.Context) error {
	return nil
}

func (ms *MemStorage) LoadFromFile(ctx context.Context) error {
	return nil
}

func (ms *MemStorage) SaveToFile(ctx context.Context) error {
	return nil
}

func (ms *MemStorage) WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return nil
}
