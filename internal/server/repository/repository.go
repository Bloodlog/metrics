package repository

import (
	"errors"
	"sync"
)

type MetricType string

type MetricStorage interface {
	SetGauge(name string, value float64) error
	GetGauge(name string) (float64, error)
	SetCounter(name string, value uint64) error
	GetCounter(name string) (uint64, error)
	Gauges() map[string]float64
	Counters() map[string]uint64
	AutoSave() error
	LoadFromFile() error
	SaveToFile() error
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

func (ms *MemStorage) SetGauge(name string, value float64) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.gauges[name] = value

	return nil
}

func (ms *MemStorage) GetGauge(name string) (float64, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	value, exists := ms.gauges[name]
	if !exists {
		return 0, errors.New("gauge metric not found")
	}
	return value, nil
}

func (ms *MemStorage) SetCounter(name string, value uint64) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.counters[name] += value

	return nil
}

func (ms *MemStorage) GetCounter(name string) (uint64, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	value, exists := ms.counters[name]
	if !exists {
		return 0, errors.New("counter metric not found")
	}
	return value, nil
}

func (ms *MemStorage) Gauges() map[string]float64 {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	result := make(map[string]float64, len(ms.gauges))
	for k, v := range ms.gauges {
		result[k] = v
	}
	return result
}

func (ms *MemStorage) Counters() map[string]uint64 {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	result := make(map[string]uint64, len(ms.counters))
	for k, v := range ms.counters {
		result[k] = v
	}
	return result
}

func (ms *MemStorage) AutoSave() error {
	return nil
}

func (ms *MemStorage) LoadFromFile() error {
	return nil
}

func (ms *MemStorage) SaveToFile() error {
	return nil
}
