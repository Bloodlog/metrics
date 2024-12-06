package repository

import (
	"errors"
	"sync"
)

var (
	ErrMetricNotFound  = errors.New("gauge metric not found")
	ErrCounterNotFound = errors.New("counter metric not found")
)

type MetricType string

type MetricStorage interface {
	SetGauge(name string, value float64)
	GetGauge(name string) (float64, error)
	SetCounter(name string, value uint64)
	GetCounter(name string) (uint64, error)
}

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

func (ms *MemStorage) SetGauge(name string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.gauges[name] = value
}

func (ms *MemStorage) GetGauge(name string) (float64, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	value, exists := ms.gauges[name]
	if !exists {
		return 0, ErrMetricNotFound
	}
	return value, nil
}

func (ms *MemStorage) SetCounter(name string, value uint64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.counters[name] += value
}

func (ms *MemStorage) GetCounter(name string) (uint64, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	value, exists := ms.counters[name]
	if !exists {
		return 0, ErrCounterNotFound
	}
	return value, nil
}

func (ms *MemStorage) Gauges() map[string]float64 {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.gauges
}

func (ms *MemStorage) Counters() map[string]uint64 {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.counters
}
